package translate

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	intFormat   = "%s=%d %d\n"
	floatFormat = "%s=%f %d\n"
)

type Translate interface {
	WSToInfluxDB(data []byte) (string, error)
}

type Translator struct{}

type WSData struct {
	SysTemperature float64 `json:"sys_temperature"`
	Temperature    float64 `json:"temperature"`
	Humidity       float64 `json:"humidity"`
	Pressure       float64 `json:"pressure"`
	Light          float64 `json:"light"`
	WindSpeed      float64 `json:"wind_speed"`
	WindDirection  int     `json:"wind_direction"`
	RainTotal      float64 `json:"rain_total"`
	Rain           float64 `json:"rain"`
}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) WSToInfluxDB(data []byte) (string, error) {
	wsData := WSData{}
	if err := json.Unmarshal(data, &wsData); err != nil {
		return "", err
	}
	influxDBData := fmt.Sprintf(floatFormat, "sys_temperature", wsData.SysTemperature, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "temperature", wsData.Temperature, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "humidity", wsData.Humidity, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "pressure", wsData.Pressure, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "light", wsData.Light, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "wind_speed", wsData.WindSpeed, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(intFormat, "wind_direction", wsData.WindDirection, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "rain_total", wsData.RainTotal, time.Now().UnixNano())
	influxDBData += fmt.Sprintf(floatFormat, "rain", wsData.Rain, time.Now().UnixNano())
	return influxDBData, nil
}

// TODO: Time should be sent from collector
// TODO: consts for key names
// TODO: fixed format might be better as string template
