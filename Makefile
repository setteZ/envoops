
deploy:
	@bash -c "source ./utils/.venv/bin/activate && bash ./utils/deploy.sh"

repl:
	picocom /dev/ttyUSB0 -b115200
