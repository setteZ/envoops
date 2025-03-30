#!/bin/bash

if [ $# != 1 ]; then
  echo "missing argument"
  exit 1
else
  if [ -f "$1" ]; then
    source .venv/bin/activate
    esptool.py --chip esp32 --port /dev/ttyUSB0 erase_flash
    esptool.py --chip esp32 --port /dev/ttyUSB0 --baud 460800 write_flash 0x1000 $1
    deactivate
  else
    echo "the file does not exist"
  fi
fi
