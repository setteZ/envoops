from machine import SDCard

CONFIG_PATH = "/sd/config.json"
data = {
  "host": "",
  "network": {
    "ssid": "",
    "pwd": "",
    "addr4": "",
    "gw4": ""
  },
  "mqtt": {
    "server": "",
    "port": "",
    "topic": ""
  },
  "ota_update": {
    "host": "",
    "prj": ""
  }
}
sd = SDCard(slot=2)
