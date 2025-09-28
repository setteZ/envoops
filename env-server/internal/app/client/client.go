package client

import (
	"encoding/json"
	"env-server/internal/database"
	"env-server/models"
	"env-server/utils"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var DEFAULT_BROKER_ADDRESS = "127.0.0.1"
var DEFAULT_BROKER_PORT = "1883"

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	log.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())

	type MQTT_In_Data struct {
		Value    float64 `json:"value"`
		Quantity string  `json:"quantity"`
	}

	var mqtt_data MQTT_In_Data
	json.Unmarshal([]byte(message.Payload()), &mqtt_data)

	var data models.NodeData
	data.NodeId = strings.Split(message.Topic(), "/")[1]
	data.Quantity = mqtt_data.Quantity
	data.Value = mqtt_data.Value
	data.Time = time.Now().Format("2006-01-02_15:04:05")
	database.AddData(&data)
}

func Run() {
	hostname, _ := os.Hostname()
	topic := flag.String("topic", "node/+", "The topic name to/from which to publish/subscribe")
	broker_uri := utils.GetEnv("ENVBROKER_ADDRESS", DEFAULT_BROKER_ADDRESS) + ":" + utils.GetEnv("ENVBROKER_PORT", DEFAULT_BROKER_PORT)
	broker := flag.String("broker", broker_uri, "The broker URI")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID(*clientid)
	opts.SetCleanSession(true)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
