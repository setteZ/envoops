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

import sht30

i2c = I2C(1, sda=Pin(21), scl=Pin(22))
sht = sht30.SHT30(i2c=i2c, i2c_address=68)
led = Pin(2, Pin.OUT)


def do_connect(ssid: str, pwd: str, addr4="", gw4=""):
    """connect to network"""
    sta_if = network.WLAN(network.WLAN.IF_STA)
    addr_p = re.compile("\\d+\\.\\d+\\.\\d+\\.\\d+/\\d+")
    gw_p = re.compile("\\d+\\.\\d+\\.\\d+\\.\\d+")
    if addr_p.match(addr4) and gw_p.match(gw4):
        sta_if.ipconfig(addr4=addr4, gw4=gw4)
    if not sta_if.isconnected():
        print("connecting to network...")
        sta_if.active(True)
        sta_if.connect(ssid, pwd)
        while not sta_if.isconnected():
            machine.idle()
    print("network config:", sta_if.ipconfig("addr4"))


def web_page(host: str):
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
    t_int, t_dec, h_int, h_dec = sht.measure_int()
    t = f"{t_int}.{t_dec:02d}"
    h = f"{h_int}.{h_dec:02d}"
    resp = html % (host, t, h)
    return resp


# read wifi credential from SD card
sd = SDCard(slot=2)
led.on()
try:
    os.mount(sd, "/sd")  # mount
except:  # pylint: disable=bare-except
    print("missing uSD")
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()
try:
    with open("/sd/config.json", "r", encoding="utf-8") as f:
        data = json.load(f)
except:  # pylint: disable=bare-except
    os.umount("/sd")  # eject
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()

os.umount("/sd")  # eject
print(data)

try:
    SSID = data["network"]["ssid"]
    PWD = data["network"]["pwd"]
except:  # pylint: disable=bare-except
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()
try:
    HOST = data["host"]
except:
    HOST = "envoops node"
try:
    ADDR4 = data["network"]["addr4"]
    GW4 = data["network"]["gw4"]
except:
    ADDR4 = ""
    GW4 = ""
# connect to the wifi
led.off()
time.sleep_ms(200)
led.on()
do_connect(SSID, PWD, ADDR4, GW4)
led.off()

# create the webserver
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind(("", 80))
s.listen(5)

while True:
    conn, addr = s.accept()
    led.on()
    print(f"Got a connection from {addr}")
    request = conn.recv(1024)
    response = web_page(HOST)
    conn.send("HTTP/1.1 200 OK\n")
    conn.send("Content-Type: text/html\n")
    conn.send("Connection: close\n\n")
    conn.sendall(response)
    conn.close()
    led.off()
