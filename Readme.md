# envoops

The purpose of this repo is to develop an environmental wifi node using and ESP32 based board and [micropython](https://micropython.org/) for rapid prototyping.

## MicroPython

### Prerequisites

The following assume that `/dev/ttyUSB0` is the seial port of the board.
- install [esptool](https://docs.espressif.com/projects/esptool/en/latest/esp32/)
- install [ampy](https://pypi.org/project/adafruit-ampy/)
- allow user to access `/dev/ttyUSB0` with `sudo usermod -a -G dialout $user` or the right group
- verify that everything is fine with `esptool.py flash_id` (or `esptool.py -p /dev/ttyUSB0 flash_id` esptool.py can't automatically detect the serial port)
- install `picocom` for REPL access

### Installation

Follow the instruction [here](https://micropython.org/download/ESP32_GENERIC/) to flash the board.
