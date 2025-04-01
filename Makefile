
deploy:
	@bash -c "source ./utils/.venv/bin/activate && bash ./utils/deploy.sh"

repl:
	picocom /dev/ttyUSB0 -b115200

run-broker:
	@bash -c "mosquitto -c ./broker/mosquitto.conf"
