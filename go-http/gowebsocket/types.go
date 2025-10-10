package gowebsocket

import (
	"time"
)

// 消息类型常量
// 消息类型常量
const (
	MsgTypeHeartbeat    = "heartbeat"
	MsgTypeHeartbeatAck = "heartbeat_ack"
	MsgTypeChat         = "chat"
	MsgTypeSystem       = "system"
	MsgTypeShutdown     = "shutdown"
)

// 消息结构体
type Message struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	ClientID  string    `json:"client_id"`
	Timestamp time.Time `json:"timestamp"`
}

// 客户端连接结构
