# env-node

An environmental wifi node consisting of an ESP32 based board plus a temperature&humidity sensor  developend in [micropython](https://micropython.org/) for rapid prototyping.

## Third-Party Code

This project includes third-party components included via Git submodule:

- [micropython_ota](https://github.com/settez/micropython_ota/tree/8420316b833a871c912032a0d45d295cb53c7d40) (MIT License) - submodule at `env-node/lib/micropython_ota/`.
- [micropython-mqtt](https://github.com/chrismoorhouse/micropython-mqtt/tree/df542c8bedcb4daf98239813e6f424d90ccdae78) (BSD 3-Clause License) - submodule at `env-node/lib/micropython-mqtt/`.
- [sht30](https://github.com/robert-hh/SHT30/tree/0352fe9513fcc96a7bfaba8edb0cccccd2d8b0f8) (Apache 2.0 License) - submodule at `env-node/lib/sht30/`.

See each submodule's `LICENSE` file for details.

## Prerequisites

The following assume that `/dev/ttyUSB0` is the serial port of the board.
- install [esptool](https://docs.espressif.com/projects/esptool/en/latest/esp32/)
- install [ampy](https://pypi.org/project/adafruit-ampy/)
- allow user to access `/dev/ttyUSB0` with `sudo usermod -a -G dialout $user` or the right group
- verify that everything is fine with `esptool.py flash_id` (or `esptool.py -p /dev/ttyUSB0 flash_id` if esptool.py can't automatically detect the serial port)
- install `picocom` for REPL access

## Installation

Follow the instruction [here](https://micropython.org/download/ESP32_GENERIC/) to flash the board.
