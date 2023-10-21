import time
import paho.mqtt.client as mqtt
import json
import weatherhat

# MQTT setup
MQTT_HOST = "10.0.0.41"
MQTT_PORT = 1883
MQTT_KEEPALIVE_INTERVAL = 45
MQTT_TOPIC = "v1/weatherstation/data"

# Define on_publish event function
def on_publish(client, userdata, mid):
    print("Message Published...")

# Initiate MQTT Client
mqttc = mqtt.Client()
# Connect with MQTT Broker
mqttc.connect(MQTT_HOST, MQTT_PORT, MQTT_KEEPALIVE_INTERVAL)

# Configure weatherhat
sensor = weatherhat.WeatherHAT()

while True:
    sensor.update(interval=5.0)
    # Read sensor data
    data = {
        "sys_temperature": sensor.device_temperature, # celsius
        "temperature": sensor.temperature, # celsius
        "humidity": sensor.humidity, # %
        "pressure": sensor.pressure, # hPa
        "light": sensor.lux,
        "wind_speed": sensor.wind_speed, # mph
        "wind_direction": sensor.wind_direction, # degrees
        "rain_tota": sensor.rain_total, # mm
        "rain": sensor.rain, # mm/sec
    }
    mqttc.publish(MQTT_TOPIC, json.dumps(data))
    time.sleep(5.0)