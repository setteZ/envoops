"""
This is the main
"""

import asyncio
import time

import machine
from machine import Pin, SDCard
import network

from mqtt import MQTTClient

import micropython_ota

import configs
import server
from utils import get_version

led = Pin(2, Pin.OUT)

current_version = get_version()

def error():
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()

HOST = configs.data["host"]
if HOST == "":
    global mac
    HOST = f"env-node_{mac}"

MQTT_SERVER = configs.data["mqtt"]["server"]
MQTT_PORT = configs.data["mqtt"]["port"]
if MQTT_PORT == "":
    MQTT_PORT = 1883

MQTT_PUBLISH_TOPIC = configs.data["mqtt"]["topic"]
if MQTT_PUBLISH_TOPIC == "":
    MQTT_PUBLISH_TOPIC = HOST

OTA_HOST = configs.data["ota_update"]["host"]
OTA_PRJ = configs.data["ota_update"]["prj"]


def puback_cb(msg_id):
  print('PUBACK ID = %r' % msg_id)

def suback_cb(msg_id, qos):
  print('SUBACK ID = %r, Accepted QOS = %r' % (msg_id, qos))
  
def con_cb(connected):
  if connected:
    client.subscribe('cmd/+')

def msg_cb(topic, pay):
  topic_str = topic.decode("utf-8")
  pay_str = pay.decode("utf-8")
  print('Received %s: %s' % (topic_str, pay_str))
  dest = topic_str.split("/")[1]
  if dest in ["all", HOST]:
    print(f"the command {pay_str} is for me")
    if pay_str == "update":
        if (OTA_PRJ != "") and (OTA_HOST != ""):
            micropython_ota.check_for_ota_update(OTA_HOST, OTA_PRJ)

if MQTT_SERVER != "":
    client = MQTTClient(MQTT_SERVER, port=MQTT_PORT)

    client.set_connected_callback(con_cb)
    client.set_puback_callback(puback_cb)
    client.set_suback_callback(suback_cb)
    client.set_message_callback(msg_cb)

    client.connect(HOST)

    timeout = 100  # ~10 seconds
    while timeout > 0 and not client.isconnected():
        time.sleep_ms(100)
        timeout -= 1

asyncio.run(server.main())
