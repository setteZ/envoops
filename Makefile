SHELL=/bin/bash
VERSION := $(shell cat ./env-node/version | tr -d '\n')

put:
	@bash -c "source ./utils/.venv/bin/activate && bash ./utils/put.sh"

repl:
	picocom /dev/ttyUSB0 -b115200

run-broker:
	@bash -c "mosquitto -c ./broker/mosquitto.conf"

.PHONY: release
release: clean
	mkdir -p release-node/$(VERSION)
	python ./utils/make-release.py -b ./utils/put.sh -n -o release-node/$(VERSION)

deploy-ota: release
	@if [ -f .env ]; then \
		echo "Sourcing .env file..."; \
		set -a; \
		. .env; \
		set +a; \
	fi; \
	./utils/deploy_node_update.sh $(VERSION)
	@echo "done"

clean:
	@rm -r release-node
