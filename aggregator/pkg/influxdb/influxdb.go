package influxdb

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	pingAPI  = "/ping"
	writeAPI = "/write"
	queryAPI = "/query"
	dbName   = "weather"
)

type AGInfluxDB interface {
	CheckConnection() error
	WriteData(data []byte) error
	CreateDB(name string) error
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
		return fmt.Errorf("InfluxDB write error - %s", resp.Status())
	}

	return nil
}

func (c *AGInfluxDBClient) CreateDB(name string) error {
	params := make(map[string]string)
	params["q"] = fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetQueryParams(params).
		Post(queryAPI)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("InfluxDB write error - %s", resp.Status())
	}
	return nil
}
