#! /bin/bash
ampy --port /dev/ttyUSB0 put ./env-node/main.py
ampy --port /dev/ttyUSB0 put ./env-node/sht30/sht30.py
ampy --port /dev/ttyUSB0 put ./env-node/micropython-mqtt/src/mqtt.py
