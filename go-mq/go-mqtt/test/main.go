package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	gomqtt "github.com/gif-gif/go.io/go-mq/go-mqtt"
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
		Server:           "182.92.216.214",
		Port:             18830,
		ClientId:         "go_mqtt_client_test",
		User:             "mqtt",
		Secret:           "223238",
		DefaultHandler:   &messagePubHandler,
		OnConnect:        &connectHandler,
		OnConnectionLost: &connectLostHandler,
	}

	err := gomqtt.Init(conf)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	err = gomqtt.GetClient().Subscribe("topic/test", 1, func(c mqtt.Client, msg mqtt.Message) {
		golog.WithTag("mqtt").Info(string(msg.Payload()))
	})
	if err != nil {
		golog.Error(err.Error())
		return
	}
	//time.Sleep(5 * time.Second)
	publish(gomqtt.GetClient())
	//gomqtt.GetClient().Disconnect(250)
	<-gocontext.Cancel().Done()
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
		//time.Sleep(time.Second)
		golog.WithTag("mqtt-publish").Info(text)
	}
}
