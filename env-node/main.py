"""
This is the main
"""

import json
import os
import re
import socket
import time

import machine
from machine import I2C, Pin, SDCard
import network

from mqtt import MQTTClient
import sht30

i2c = I2C(1, sda=Pin(21), scl=Pin(22))
sht = sht30.SHT30(i2c=i2c, i2c_address=68)
led = Pin(2, Pin.OUT)

def error():
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()


def web_page(host: str, t: str, h: str):
    """prepare response page"""
    html = """<!DOCTYPE html>
<html>
    <head> <title>envoops</title> </head>
    <body> 
        <h2>%s</h2>
        <h3>STH3X</h3>
        <p>temperature: %s degC</p>
        <p>humidity: %s %% RH</p>
    </body>
</html>
"""
    resp = html % (host, t, h)
    return resp


# read wifi credential from SD card
os.mount(sd, "/sd")  # mount
with open("/sd/config.json", "r", encoding="utf-8") as f:
    data = json.load(f)
os.umount("/sd")  # eject
print(data)

try:
    HOST = data["host"]
except:
    global mac
    HOST = f"env-node_{mac}"
try:
    MQTT_SERVER = data["mqtt"]["server"]
    MQTT_PORT = data["mqtt"]["port"]
except:
    print("missing mqtt server and port info")
    error()
try:
    MQTT_PUBLISH_TOPIC = data["mqtt"]["topic"]
except:
    MQTT_PUBLISH_TOPIC = HOST


def puback_cb(msg_id):
  print('PUBACK ID = %r' % msg_id)

def suback_cb(msg_id, qos):
  print('SUBACK ID = %r, Accepted QOS = %r' % (msg_id, qos))
  
def con_cb(connected):
  if connected:
    client.subscribe('subscribe/topic')

def msg_cb(topic, pay):
  print('Received %s: %s' % (topic.decode("utf-8"), pay.decode("utf-8")))

# create the webserver
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind(("", 80))
s.listen(5)

client = MQTTClient(MQTT_SERVER, port=MQTT_PORT)

client.set_connected_callback(con_cb)
client.set_puback_callback(puback_cb)
client.set_suback_callback(suback_cb)
client.set_message_callback(msg_cb)

client.connect(HOST)

while not client.isconnected():
    time.sleep_ms(100)

while True:
    conn, addr = s.accept()
    led.on()
    print(f"Got a connection from {addr}")
    request = conn.recv(1024)
    request_str = request.decode('utf-8')
    if request_str.find('/info') > 0:
        conn.send("HTTP/1.1 200 OK\n")
        conn.send("Content-Type: application/json\n")
        conn.send("Connection: close\n\n")
        conn.sendall(json.dumps(data))
        conn.close()
    else:
        t_int, t_dec, h_int, h_dec = sht.measure_int()
        t = f"{t_int}.{t_dec:02d}"
        h = f"{h_int}.{h_dec:02d}"
        response = web_page(HOST, t, h)
        conn.send("HTTP/1.1 200 OK\n")
        conn.send("Content-Type: text/html\n")
        conn.send("Connection: close\n\n")
        conn.sendall(response)
        conn.close()
        if client.isconnected():
            try:
                pub_id = client.publish(f"{MQTT_PUBLISH_TOPIC}/temperature", t, False)
                pub_id = client.publish(f"{MQTT_PUBLISH_TOPIC}/humidity", h, False)
            except Exception as e:
                print(f"exception: {e}")
    led.off()
