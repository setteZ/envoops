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
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var DEFAULT_BROKER_ADDRESS = "127.0.0.1"
var DEFAULT_BROKER_PORT = "1883"

func Run() {
	topic := flag.String("topic", "node/+", "The topic name to/from which to publish/subscribe")
	broker_uri := utils.GetEnv("ENVBROKER_ADDRESS", DEFAULT_BROKER_ADDRESS) + ":" + utils.GetEnv("ENVBROKER_PORT", DEFAULT_BROKER_PORT)
	broker := flag.String("broker", broker_uri, "The broker URI")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	receiveCount := 0
	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for receiveCount < *num {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])

		type MQTT_In_Data struct {
			Value    float64 `json:"value"`
			Quantity string  `json:"quantity"`
		}

		var mqtt_data MQTT_In_Data
		json.Unmarshal([]byte(incoming[1]), &mqtt_data)

		var data models.NodeData
		data.Quantity = mqtt_data.Quantity
		data.Value = mqtt_data.Value
		data.Time = time.Now().Format("2006-01-02_15:04:05")
		log.Print(incoming[1])
		log.Print(data)
		database.AddData(&strings.Split(incoming[0], "/")[1], &data)
		receiveCount++
	}

	client.Disconnect(250)
	fmt.Println("Sample Subscriber Disconnected")
}
