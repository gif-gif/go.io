package golog

type Adapter interface {
	Write(msg *Message)
}
