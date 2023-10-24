package influxdb

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	pingAPI  = "/ping"
	writeAPI = "/write"
	dbName   = "weather"
)

type AGInfluxDB interface {
	CheckConnection() error
	WriteData(data []byte) error
}

type AGInfluxDBClient struct {
	client *resty.Client
}

func NewAGInfluxDBClient(username, password, url string) *AGInfluxDBClient {
	restClient := resty.New()
	restClient.SetBasicAuth(username, password)
	restClient.SetTimeout(time.Second * 5)
	restClient.SetBaseURL(url)
	return &AGInfluxDBClient{client: restClient}
}

func (c *AGInfluxDBClient) CheckConnection() error {
	resp, err := c.client.R().
		SetHeader("Content-Type", "none").
		SetHeader("Accept", "application/json").
		Get(pingAPI)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return err
	}
	return nil
}

func (c *AGInfluxDBClient) WriteData(data []byte) error {
	var params = map[string]string{
		"db":        dbName,
		"precision": "s",
	}
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(data).
		SetQueryParams(params).
		Post(writeAPI)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return err
	}
	return nil
}
