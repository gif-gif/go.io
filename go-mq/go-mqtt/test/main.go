package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	gocontext "github.com/gif-gif/go.io/go-context"
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

var conf = gomqtt.Config{
	Server:           "182.92.216.214",
	Port:             18830,
	ClientId:         "go_mqtt_client_test",
	User:             "mqtt",
	Secret:           "223238",
	DefaultHandler:   &messagePubHandler,
	OnConnect:        &connectHandler,
	OnConnectionLost: &connectLostHandler,
}

var topic = "topic/test"

func main() {
	testSubscribe()
	time.Sleep(time.Second * 2)
	testPublish()
	<-gocontext.Cancel().Done()
}

func testSubscribe() {
	conf.Name = "testSubscribe"
	err := gomqtt.Init(conf)
	if err != nil {
		golog.Error(err.Error())
		return
	}
	c := gomqtt.GetClient(conf.Name)
	c.Client.Unsubscribe(topic)
	err = c.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
		golog.WithTag("mqtt-Subscribe").Info(string(msg.Payload()))
	})
	if err != nil {
		golog.Error(err.Error())
		return
	}
}

func testPublish() {
	conf.Name = "testPublish"
	err := gomqtt.Init(conf)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	publish(gomqtt.GetClient(conf.Name))
	<-gocontext.Cancel().Done()
}

func publish(client *gomqtt.GoMqttClient) {
	num := 10

	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message1111 %d", i)
		err := client.Publish(topic, 0, true, text)
		if err != nil {
			golog.Error(err.Error())
			return
		}
		time.Sleep(time.Second)
		golog.WithTag("mqtt-publish").WithField(topic, topic).Info(text)
	}
	err := client.Publish(topic, 0, true, "publish-text")
	if err != nil {
		golog.Error(err.Error())
		return
	}
	golog.WithTag("mqtt-publish").WithField(topic, topic).Info("publish-text")
}
