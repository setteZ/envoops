#! /bin/bash
ampy --port /dev/ttyUSB0 put ./env-node/boot.py
ampy --port /dev/ttyUSB0 put ./env-node/main.py
ampy --port /dev/ttyUSB0 put ./env-node/utils.py
ampy --port /dev/ttyUSB0 put ./env-node/configs.py
ampy --port /dev/ttyUSB0 put ./env-node/server.py
ampy --port /dev/ttyUSB0 put ./env-node/measurement.py
ampy --port /dev/ttyUSB0 put ./env-node/lib/sht30/sht30.py
ampy --port /dev/ttyUSB0 put ./env-node/lib/micropython-mqtt/src/mqtt.py
ampy --port /dev/ttyUSB0 put ./env-node/version
ampy --port /dev/ttyUSB0 put ./env-node/lib/micropython_ota/micropython_ota.py
