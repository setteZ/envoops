# envoops

The purpose of this repo is to develop an environmental wifi node using and ESP32 based board and [micropython](https://micropython.org/) for rapid prototyping.

## MicroPython

### Prerequisites

The following assume that `/dev/ttyUSB0` is the seial port of the board.
- install [esptool](https://docs.espressif.com/projects/esptool/en/latest/esp32/)
- install [ampy](https://pypi.org/project/adafruit-ampy/)
- allow user to access `/dev/ttyUSB0` with `sudo usermod -a -G dialout $user` or the right group
- verify that everything is fine with `esptool.py flash_id` (or `esptool.py -p /dev/ttyUSB0 flash_id` if esptool.py can't automatically detect the serial port)
- install `picocom` for REPL access

### Installation

Follow the instruction [here](https://micropython.org/download/ESP32_GENERIC/) to flash the board.

## Utils

With the `Makefile` there is the possibility to call some scripts to make the life easier:

- `make put` update the board connected to the USB
- `make repl` open the serial repl connection (using `picocom`)
- `make run-broker` run localli a MQTT broker for manual test
- `make release` prepare a folder with all the files necessary for a release
- `make deploy-ota` deploy on a server the release files

A `.env` file (that you can find in the template folder as reference) is usefull for the configuration of the server references.
