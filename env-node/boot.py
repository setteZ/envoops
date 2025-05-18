"""
This is the main
"""
import gc
gc.mem_free()
import json
import os
import re
import time

import ubinascii

import machine
from machine import Pin, SDCard
import network

import micropython_ota

led = Pin(2, Pin.OUT)

global mac

def error():
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()

def do_connect(ssid: str, pwd: str, addr4="", gw4="", dns="0.0.0.0") -> str:
    """connect to network"""
    sta_if = network.WLAN(network.WLAN.IF_STA)
    mac = ubinascii.hexlify(sta_if.config('mac')).decode()
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
    network.ipconfig(dns=dns)
    return mac[:6]


# read wifi credential from SD card
sd = SDCard(slot=2)
led.on()
try:
    os.mount(sd, "/sd")  # mount
except:  # pylint: disable=bare-except
    print("missing uSD")
    error()
try:
    with open("/sd/config.json", "r", encoding="utf-8") as f:
        data = json.load(f)
except:  # pylint: disable=bare-except
    os.umount("/sd")  # eject
    error()

os.umount("/sd")  # eject
print(data)

try:
    SSID = data["network"]["ssid"]
    PWD = data["network"]["pwd"]
except:  # pylint: disable=bare-except
    error()
try:
    ADDR4 = data["network"]["addr4"]
    GW4 = data["network"]["gw4"]
except:
    ADDR4 = ""
    GW4 = ""

try:
    DNS = data["network"]["dns"]
except:
    DNS = GW4

try:
    OTA_HOST = data["ota_update"]["host"]
    OTA_PRJ = data["ota_update"]["prj"]
except:
    OTA_HOST = None
    OTA_PRJ = None

# connect to the wifi
led.off()
time.sleep_ms(200)
led.on()
mac = do_connect(SSID, PWD, ADDR4, GW4, DNS)
led.off()
if (OTA_PRJ != None) and (OTA_HOST != None):
    micropython_ota.ota_update(OTA_HOST, OTA_PRJ, use_version_prefix=False)
