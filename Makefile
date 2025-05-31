
VERSION := $(shell cat ./env-node/version | tr -d '\n')

put:
	@bash -c "source ./utils/.venv/bin/activate && bash ./utils/put.sh"

repl:
	picocom /dev/ttyUSB0 -b115200

run-broker:
	@bash -c "mosquitto -c ./broker/mosquitto.conf"

.PHONY: release
release:
	mkdir -p release-node/$(VERSION)
	@bash -c "python ./utils/make-manifest.py -b ./utils/put.sh -n -o release-node/$(VERSION)"
