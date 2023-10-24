package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ConnectionOptions struct {
	BrokerScheme     string
	BrokerAddress    string
	BrokerPort       int
	ClientID         string
	CleanSession     bool
	AutoReconnect    bool
	ConnectRetry     bool
	Store            mqtt.Store
	ResumeSubs       bool
	OnConnect        mqtt.OnConnectHandler
	OnConnectionLost mqtt.ConnectionLostHandler
}
