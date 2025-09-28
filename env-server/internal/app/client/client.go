package client

import (
	"encoding/json"
	"env-server/internal/database"
	"env-server/models"
	"env-server/utils"
	"flag"
	"fmt"
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

	data, err := parseMessage(message.Topic(), message.Payload())
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	database.AddData(data)
}

type MQTT_In_Data struct {
	Value    float64 `json:"value"`
	Quantity string  `json:"quantity"`
}

func parseMessage(topic string, payload []byte) (*models.NodeData, error) {
	var mqttData MQTT_In_Data
	if err := json.Unmarshal(payload, &mqttData); err != nil {
		return nil, err
	}

	parts := strings.Split(topic, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid topic format: %s", topic)
	}

	return &models.NodeData{
		NodeId:   parts[1],
		Quantity: mqttData.Quantity,
		Value:    mqttData.Value,
		Time:     time.Now().Format("2006-01-02_15:04:05"),
	}, nil
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
