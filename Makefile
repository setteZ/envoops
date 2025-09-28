SHELL=/bin/bash
NODE_VERSION := $(shell cat ./env-node/version | tr -d '\n')

node-put:
	@bash -c "source ./utils/.venv/bin/activate && bash ./utils/put.sh"

node-repl:
	picocom /dev/ttyUSB0 -b115200

.PHONY: node-release
node-release: clean
	mkdir -p node-release/$(NODE_VERSION)
	python ./utils/make-release.py -b ./utils/put.sh -n -o node-release/$(NODE_VERSION)

node-deploy-ota: release
	@if [ -f .env ]; then \
		echo "Sourcing .env file..."; \
		set -a; \
		. .env; \
		set +a; \
	fi; \
	./utils/deploy_node_update.sh $(NODE_VERSION)
	@echo "done"

node-clean:
	@rm -r -f node-release

broker-run:
	@bash -c "mosquitto -c ./broker/mosquitto.conf"
