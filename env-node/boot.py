"""
This is the main
"""
import gc
gc.mem_free()
import asyncio
import json
import os
import re
import time

import ubinascii

import machine
from machine import Pin, SDCard
import network

import micropython_ota

import configs
import server

led = Pin(2, Pin.OUT)

global mac
SSID = ""

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
else:
    try:
        with open(configs.CONFIG_PATH, "r", encoding="utf-8") as f:
            configs.data = json.load(f)
    except:  # pylint: disable=bare-except
        print(f"missing {configs.CONFIG_PATH} file")
    else:
        print(configs.data)
        try:
            PWD = configs.data["network"]["pwd"]
            SSID = configs.data["network"]["ssid"]
        except:  # pylint: disable=bare-except
            print("missing info for the SSID connection")
        try:
            ADDR4 = configs.data["network"]["addr4"]
            GW4 = configs.data["network"]["gw4"]
        except:  # pylint: disable=bare-except
            ADDR4 = ""
            GW4 = ""
        try:
            DNS = configs.data["network"]["dns"]
        except:  # pylint: disable=bare-except
            DNS = GW4
        try:
            OTA_HOST = configs.data["ota_update"]["host"]
            OTA_PRJ = configs.data["ota_update"]["prj"]
        except:  # pylint: disable=bare-except
            OTA_HOST = None
            OTA_PRJ = None

finally:
    os.umount("/sd")  # eject


# connect to the wifi
led.off()
time.sleep_ms(200)
led.on()
if SSID != "":
    mac = do_connect(SSID, PWD, ADDR4, GW4, DNS)
else:
    # Configure access point with MAC-based SSID
    ap = network.WLAN(network.WLAN.IF_AP) # create access-point interface
    mac = ubinascii.hexlify(ap.config('mac')).decode()
    print(f"{mac}")
    ap.config(ssid=f"ESP32-{mac[:6]}")    # set the SSID of the access point
    ap.config(max_clients=1)              # set how many clients can connect to the network
    ap.active(True)                       # activate the interface
    while ap.active() is False:
        pass
    asyncio.run(server.main())
    while True:
        time.sleep(10)
led.off()
if (OTA_PRJ != None) and (OTA_HOST != None):
    micropython_ota.ota_update(OTA_HOST, OTA_PRJ, use_version_prefix=False)
