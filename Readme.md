# envoops

The purpose of this project is to develop an environmental monitoring with sensor network.

## Repository structure

This is a **monorepo** divided into independent but interrelated components:
- [env-node](env-node) is the physcal node the perform the measurements (written in Micropython)
- [env-server](env-server) is the server that collect and save all the data (written in Go)
- [env-broker](env-broker) the MQTT broker that dispatch the messagges to the network (using [mosquitto](https://mosquitto.org/))

Once you clone the repo, you shall update the submodules:

```bash
gt clone https://github.com/setteZ/envoops.git
cd envoops
git submodule update --init --recursive
```

Each component has its own release cycle and versioning, and we use Git tags in the format:

```
<component-name>/v<semver>
```

Example tags:
- `env-node/v1.3.0`
- `env-server/v2.0.1`

See each component's README for specific setup and dependencies.

## Utils

With the `Makefile` there is the possibility to call some scripts to make the life easier:

- `make put-node` update the board connected to the USB
- `make repl` open the serial repl connection (using `picocom`)
- `make run-broker` run locally a MQTT broker for manual test
- `make release-node` prepare a folder with all the files necessary for a env-node release
- `make deploy-ota` deploy on a server the env-node release

A `.env` file (you can find a template [here](templates/.env) as reference) is usefull for the configuration of the server references.
