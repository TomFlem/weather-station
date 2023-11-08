package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TomFlem/weather-station/aggregator/internal/translate"
	"github.com/TomFlem/weather-station/aggregator/pkg/influxdb"
	amqtt "github.com/TomFlem/weather-station/aggregator/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	dataTopic = "v1/weatherstation/data"
)

func cleanup() {
}

func main() {
	fmt.Println("Starting Aggregator")

	mqttClient := amqtt.NewAGMQTTClient()
	opts := amqtt.ConnectionOptions{
		BrokerScheme:  "tcp",
		BrokerAddress: "localhost",
		BrokerPort:    1883,
		ClientID:      "aggregator",
		CleanSession:  true,
		AutoReconnect: true,
		ConnectRetry:  true,
		Store:         nil,
		ResumeSubs:    true,
		OnConnect:     nil,
		OnConnectionLost: func(client mqtt.Client, err error) {
			fmt.Println("Connection lost")
		},
	}
	mqttConnErr := mqttClient.Connect(&opts)
	if mqttConnErr != nil {
		fmt.Println("Error connecting to MQTT broker")
		os.Exit(1)
	}

	influxdbClient := influxdb.NewAGInfluxDBClient("admin", "admin", "http://localhost:8086")
	influxdbConnErr := influxdbClient.CheckConnection()
	if influxdbConnErr != nil {
		fmt.Printf("Error connecting to InfluxDB - %s\n", influxdbConnErr)
		os.Exit(1)
	}
	createDBErr := influxdbClient.CreateDB("weather")
	if createDBErr != nil {
		fmt.Printf("Error creating InfluxDB database - %s\n", createDBErr)
		os.Exit(1)
	}

	translator := translate.NewTranslator()

	mqttClient.Subscribe(dataTopic, func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		iData, transERr := translator.WSToInfluxDB(msg.Payload())
		if transERr != nil {
			fmt.Printf("Error translating data - %s\n", transERr)
		}
		influxdbErr := influxdbClient.WriteData([]byte(iData))
		if influxdbErr != nil {
			fmt.Printf("Error writing to InfluxDB - %s\n", influxdbErr)
		}
	})

	// cleanup handler to quit when told to do so
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-ch
		cleanup()
		os.Exit(1)
	}()

	for {
		time.Sleep(time.Second * 5)
	}
}
