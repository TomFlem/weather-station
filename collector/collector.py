import time
import paho.mqtt.client as mqtt
import json
import weatherhat

# MQTT setup

# Load Config
config = json.load(open("config.json"))
# Start
print("Weather Station Data Collector - Starting")
print ("MQTT Broker: " + config['host'] + ":" + str(config['port']))

# Define on_publish event function
def on_publish(client, userdata, mid):
    print("Weather data published on")
def on_disconnect(client, userdata, rc):
    if rc != 0:
        print("Unexpected MQTT disconnection. Will auto-reconnect")
def on_connect(client,userdata,flags,reasonCode,properties):
    print("MQTT Client Connected")

# Initiate MQTT Client
print("Initiating MQTT Client")
mqttc = mqtt.Client("weather-collector", clean_session=True)
mqttc.on_publish = on_publish
mqttc.reconnect = 1
mqttc.on_disconnect = on_disconnect
mqttc.on_connect = on_connect
# Connect with MQTT Broker
print("Connecting to MQTT Broker...")
mqttc.connect(config['host'], config['port'], config['keepalive'])
mqttc.loop_start()
# Configure weatherhat
print("Configuring WeatherHAT")
sensor = weatherhat.WeatherHAT()

print("Startup Complete - Running")
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