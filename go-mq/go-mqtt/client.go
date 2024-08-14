package gomqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const (
	PROTOCOL_TCP = "tcp"
	PROTOCOL_SSL = "ssl"
	PROTOCOL_WS  = "ws"
	PROTOCOL_WSS = "wss"
)

type Config struct {
	Name             string `yaml:"Name,optional" json:"name,optional"`
	Protocol         string `yaml:"Protocol,optional" json:"protocol,optional"`
	CaFilePath       string `yaml:"CaFilePath,optional" json:"caFilePath,optional"` //tls(ssl wss) 证书位置
	Server           string `yaml:"Server,optional" json:"server,optional"`
	Port             int    `yaml:"Port,optional" json:"port,optional"`
	ClientId         string `yaml:"ClientId,optional" json:"clientId,optional"`
	User             string `yaml:"User,optional" json:"user,optional"`
	Secret           string `yaml:"Secret,optional" json:"secret,optional"`
	KeepAlive        int64  `yaml:"KeepAlive,optional" json:"keepAlive,optional"` //s
	AutoReconnect    bool   `yaml:"AutoReconnect,optional" json:"autoReconnect,optional"`
	ConnectRetry     bool   `yaml:"ConnectRetry,optional" json:"connectRetry,optional"`
	ConnectTimeout   int64  `yaml:"ConnectTimeout,optional" json:"connectTimeout,optional"`
	DefaultHandler   *mqtt.MessageHandler
	OnConnect        *mqtt.OnConnectHandler
	OnConnectionLost *mqtt.ConnectionLostHandler
}

type GoMqttClient struct {
	Client mqtt.Client
}

func NewClient(config Config) (*GoMqttClient, error) {
	if config.ConnectTimeout <= 0 {
		config.ConnectTimeout = 10
	}

	if config.KeepAlive <= 0 {
		config.KeepAlive = 60 * 2
	}

	if config.Protocol == "" {
		config.Protocol = "tcp"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%d/mqtt", config.Protocol, config.Server, config.Port))

	if config.ClientId == "" {
		rand.Seed(time.Now().UnixNano())
		config.ClientId = fmt.Sprintf("go-client-%d", rand.Int())
	}

	opts.SetClientID(config.ClientId)
	opts.SetUsername(config.User)
	opts.SetPassword(config.Secret)
	opts.SetAutoReconnect(config.AutoReconnect)
	opts.SetConnectRetry(config.ConnectRetry)
	opts.SetKeepAlive(time.Duration(config.KeepAlive) * time.Second)
	opts.SetConnectTimeout(time.Duration(config.ConnectTimeout) * time.Second)

	// Optional: 设置CA证书
	if config.CaFilePath != "" {
		opts.SetTLSConfig(loadTLSConfig(config.CaFilePath))
	}

	if config.DefaultHandler != nil {
		opts.SetDefaultPublishHandler(*config.DefaultHandler)
	}

	if config.OnConnectionLost != nil {
		opts.SetConnectionLostHandler(*config.OnConnectionLost)
	}
	if config.OnConnect != nil {
		opts.SetOnConnectHandler(*config.OnConnect)
	}

	gc, err := NewConfig(opts)
	if err != nil {
		return nil, err
	}
	return gc, nil
}

func NewConfig(opts *mqtt.ClientOptions) (*GoMqttClient, error) {
	gomqtt := &GoMqttClient{}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	gomqtt.Client = client
	return gomqtt, nil
}

func loadTLSConfig(caFile string) *tls.Config {
	// load tls config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	if caFile != "" {
		certpool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
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
