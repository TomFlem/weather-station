package mqtt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type AGMQTT interface {
	Connect() error
	Disconnect()
	IsConnected() bool
	Publish(topic string, payload interface{}) error
	Subscribe(topic string, callback mqtt.MessageHandler) error
	Unsubscribe(topic string) error
}

type AGMQTTClient struct {
	client mqtt.Client
	mu     sync.RWMutex
}

func NewAGMQTTClient() *AGMQTTClient {
	return &AGMQTTClient{}
}

func (c *AGMQTTClient) Connect(connOpts *ConnectionOptions) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("Connecting to MQTT broker")

	// Setup Client options
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s://%s:%d", connOpts.BrokerScheme, connOpts.BrokerAddress, connOpts.BrokerPort)).
		SetClientID(connOpts.ClientID).
		SetCleanSession(connOpts.CleanSession).
		SetAutoReconnect(connOpts.AutoReconnect).
		SetConnectRetry(connOpts.ConnectRetry).
		SetStore(connOpts.Store).
		SetResumeSubs(connOpts.ResumeSubs).
		SetOnConnectHandler(connOpts.OnConnect).
		SetConnectionLostHandler(connOpts.OnConnectionLost)

	c.client = mqtt.NewClient(opts)
	token := c.client.Connect()
	if !token.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("MQTT connection timeout")
	}
	if token.Error() != nil {
		return fmt.Errorf("MQTT connection error: %v", token.Error())
	}
	return nil
}

func (c *AGMQTTClient) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.client.Disconnect(250)
}

func (c *AGMQTTClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client.IsConnected()
}

func (c *AGMQTTClient) Publish(topic string, payload interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Create Custom encoder to disable HTML escaping
	payloadJSON := &bytes.Buffer{}
	encoder := json.NewEncoder(payloadJSON)
	encoder.SetEscapeHTML(false)
	if encErr := encoder.Encode(payload); encErr != nil {
		return encErr
	}
	token := c.client.Publish(topic, 0, false, payloadJSON.Bytes())
	if !token.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("MQTT publish timeout")
	}
	if token.Error() != nil {
		return fmt.Errorf("MQTT publish error: %v", token.Error())
	}
	return nil
}

func (c *AGMQTTClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	token := c.client.Subscribe(topic, 0, callback)
	if !token.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("MQTT subscribe timeout")
	}
	if token.Error() != nil {
		return fmt.Errorf("MQTT subscribe error: %v", token.Error())
	}
	return nil
}

func (c *AGMQTTClient) Unsubscribe(topic string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	token := c.client.Unsubscribe(topic)
	if !token.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("MQTT unsubscribe timeout")
	}
	if token.Error() != nil {
		return fmt.Errorf("MQTT unsubscribe error: %v", token.Error())
	}
	return nil
}
