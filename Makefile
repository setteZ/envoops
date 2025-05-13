
VERSION := $(shell cat ./env-node/version | tr -d '\n')

deploy:
	@bash -c "source ./utils/.venv/bin/activate && bash ./utils/deploy.sh"

repl:
	picocom /dev/ttyUSB0 -b115200

run-broker:
	@bash -c "mosquitto -c ./broker/mosquitto.conf"

.PHONY: release
release:
	mkdir -p release/$(VERSION)
	@bash -c "python ./utils/make-manifest.py -b ./utils/deploy.sh -n -o release/$(VERSION)/manifest"
