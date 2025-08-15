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
from machine import Pin
import network

import micropython_ota

import configs
import server

led = Pin(2, Pin.OUT)

global mac

def error():
    while True:
        time.sleep_ms(200)
        led.off()
        time.sleep_ms(200)
        led.on()

def do_connect(ssid: str, pwd: str, addr4="", gw4="", dns="0.0.0.0") -> str | None:
    """connect to network"""
    if ssid == "":
        return None
    sta_if = network.WLAN(network.WLAN.IF_STA)
    mac = ubinascii.hexlify(sta_if.config('mac')).decode()
    addr_p = re.compile("\\d+\\.\\d+\\.\\d+\\.\\d+/\\d+")
    gw_p = re.compile("\\d+\\.\\d+\\.\\d+\\.\\d+")
    if addr_p.match(addr4) and gw_p.match(gw4):
        sta_if.ipconfig(addr4=addr4, gw4=gw4)
    if not sta_if.isconnected():
        print("connecting to network...")
        sta_if.active(True)
        sta_if.config(reconnects=1)
        sta_if.connect(ssid, pwd)
        while sta_if.status() == network.STAT_CONNECTING:
            machine.idle()
        if not sta_if.isconnected():
            sta_if.active(False)
            print("connection failed")
            return None
    print("network config:", sta_if.ipconfig("addr4"))
    network.ipconfig(dns=dns)
    return mac[:6]

def create_ap():
    # Configure access point with MAC-based SSID
    sta = network.WLAN(network.WLAN.IF_STA)
    sta.active(False)  # ensure STA mode is off so AP can work cleanly

    ap = network.WLAN(network.WLAN.IF_AP) # create access-point interface
    ap.active(True)                       # activate the interface

    mac = ubinascii.hexlify(ap.config('mac')).decode()
    print(f"{mac}")
    ap.config(ssid=f"ESP32-{mac[:6]}")    # set the SSID of the access point
    ap.config(max_clients=1)              # set how many clients can connect to the network
    #ap.active(True)                       # activate the interface
    while ap.active() is False:
        pass
    asyncio.run(server.main())
    error()

def import_config(src, dest) -> dict:
    try:
        for k,v in src.items():
            if k in dest:
                if isinstance(v, dict):
                    dest[k] = import_config(v, dest[k])
                else:
                    dest[k] = v
            else:
                print(f"key {k} not present in config")
        return dest
    except Exception as err:
        print(err)
        raise

# read wifi credential from SD card
led.on()
try:
    os.mount(configs.sd, "/sd")  # mount
except:  # pylint: disable=bare-except
    print("missing uSD")
else:
    try:
        with open(configs.CONFIG_PATH, "r", encoding="utf-8") as f:
            data = json.load(f)
        print(f"data read: {data}")
        configs.data = import_config(data, configs.data)
    except:  # pylint: disable=bare-except
        print(f"missing {configs.CONFIG_PATH} file")
    else:
        print(configs.data)

finally:
    os.umount("/sd")  # eject


# connect to the wifi
led.off()
time.sleep_ms(200)
led.on()

PWD = configs.data["network"]["pwd"]
SSID = configs.data["network"]["ssid"]
ADDR4 = configs.data["network"]["addr4"]
GW4 = configs.data["network"]["gw4"]
mac = do_connect(SSID, PWD, ADDR4, GW4, GW4)
if mac is None:
    create_ap()
led.off()
OTA_HOST = configs.data["ota_update"]["host"]
OTA_PRJ = configs.data["ota_update"]["prj"]
if (OTA_PRJ != "") and (OTA_HOST != ""):
    micropython_ota.ota_update(OTA_HOST, OTA_PRJ, use_version_prefix=False)
