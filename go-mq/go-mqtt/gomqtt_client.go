package gomqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Config struct {
	Server         string `yaml:"Server" json:"server"`
	ClientId       string `yaml:"ClientId" json:"clientId"`
	User           string `yaml:"User" json:"user"`
	Secret         string `yaml:"Secret" json:"secret"`
	KeepAlive      int64  `yaml:"KeepAlive" json:"keepAlive"` //s
	AutoReconnect  bool   `yaml:"AutoReconnect" json:"autoReconnect"`
	ConnectRetry   bool   `yaml:"ConnectRetry" json:"connectRetry"`
	ConnectTimeout int64  `yaml:"ConnectTimeout" json:"connectTimeout"`
	defaultHandler *mqtt.MessageHandler
}

type GoMqttClient struct {
	Client mqtt.Client
}

func New(config Config) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Server)
	opts.SetClientID(config.ClientId)
	opts.SetUsername(config.User)
	opts.SetPassword(config.Secret)
	opts.SetAutoReconnect(config.AutoReconnect)
	opts.SetConnectRetry(config.ConnectRetry)
	opts.SetKeepAlive(time.Duration(config.KeepAlive) * time.Second)
	opts.SetConnectTimeout(time.Duration(config.ConnectTimeout) * time.Second)
	if config.defaultHandler != nil {
		opts.SetDefaultPublishHandler(*config.defaultHandler)
	}

	err := NewConfig(opts)
	if err != nil {
		return err
	}
	return nil
}

func NewConfig(opts *mqtt.ClientOptions) error {
	gomqtt := GoMqttClient{}
	c := mqtt.NewClient(opts)
	gomqtt.Client = c
	return nil
}

func (g *GoMqttClient) Connect() error {
	if token := g.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (g *GoMqttClient) Disconnect(quiesce uint) {
	g.Client.Disconnect(quiesce)
}

func (g *GoMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := g.Client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (g *GoMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if token := g.Client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (g *GoMqttClient) Unsubscribe(topics ...string) error {
	if token := g.Client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (g *GoMqttClient) isConnected() bool {
	return g.Client.IsConnected()
}
