package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	golog "github.com/gif-gif/go.io/go-log"
	gomqtt "github.com/gif-gif/go.io/go-mq/go-mqtt"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	conf := gomqtt.Config{
		Server:           "122.228.113.238",
		Port:             1883,
		ClientId:         "go_mqtt_client_test",
		User:             "",
		Secret:           "",
		DefaultHandler:   &messagePubHandler,
		OnConnect:        &connectHandler,
		OnConnectionLost: &connectLostHandler,
	}

	gc, err := gomqtt.New(conf)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	err = gc.Subscribe("topic/test", 1, nil)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	publish(gc)
	gc.Disconnect(250)
}

func publish(client *gomqtt.GoMqttClient) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		err := client.Publish("topic/test", 0, false, text)
		if err != nil {
			golog.Error(err.Error())
			return
		}
		time.Sleep(time.Second)
	}
}
