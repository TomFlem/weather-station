package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TomFlem/weather-station/aggregator/pkg/influxdb"
	amqtt "github.com/TomFlem/weather-station/aggregator/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func cleanup() {

}
func main() {
	fmt.Println("Starting Aggregator")

	mqttClient := amqtt.NewAMQTTClient()
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

	influxdbClient := influxdb.NewAGInfluxDBClient("username", "password", "url")
	influxdbConnErr := influxdbClient.CheckConnection()
	if influxdbConnErr != nil {
		fmt.Println("Error connecting to InfluxDB")
		os.Exit(1)
	}

	// cleanup handler to quit when told to do so
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-ch
		cleanup()
		os.Exit(1)
	}()
}
