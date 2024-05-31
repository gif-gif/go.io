package goevent

type Message struct {
	Topic string
	Data  interface{}
}

type MessageChan chan Message

type SubscribeFunc func(msg Message)
