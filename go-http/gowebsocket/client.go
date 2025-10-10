package gowebsocket

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClientConfig struct {
	Name              string `yaml:"Name" json:"name,optional"`
	Port              int    `yaml:"Port" json:"port,optional"`
	Addr              string `yaml:"Addr" json:"addr,optional"`
	ClientID          string `yaml:"ClientID" json:"clientID,optional"`
	HeartBeatInterval int64  `yaml:"HeartBeatInterval" json:"heartBeatInterval,optional"`
}

// WebSocket 客户端
type WSClient struct {
	conn              *websocket.Conn
	clientID          string
	heartbeatInterval time.Duration
	lastHeartbeat     time.Time
	isConnected       bool
	mutex             sync.RWMutex
	send              chan Message // 发送通道
	done              chan struct{}
	interrupt         chan os.Signal
	ClientConfig      ClientConfig
}

// 创建新的 WebSocket 客户端
func NewWSClient(config ClientConfig) *WSClient {
	c := &config
	if c.HeartBeatInterval == 0 {
		c.HeartBeatInterval = 20
	}
	return &WSClient{
		ClientConfig:      config,
		clientID:          config.ClientID,
		heartbeatInterval: time.Duration(config.HeartBeatInterval) * time.Second, // 每20秒发送一次心跳
		isConnected:       false,
		done:              make(chan struct{}),
		send:              make(chan Message, 256),
		interrupt:         make(chan os.Signal, 1),
	}
}

// 连接到服务器
func (c *WSClient) Connect() error {
	u := url.URL{Scheme: "ws", Host: c.ClientConfig.Addr + ":" + gconv.String(c.ClientConfig.Port), Path: "/ws"}
	logx.Infof("[%s] 正在连接到 %s", c.clientID, u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}

	c.mutex.Lock()
	c.conn = conn
	c.isConnected = true
	c.lastHeartbeat = time.Now()
	c.mutex.Unlock()

	logx.Infof("[%s] 连接成功", c.clientID)

	// 设置信号处理
	signal.Notify(c.interrupt, os.Interrupt)

	return nil
}

// 启动客户端
func (c *WSClient) Start() {
	// 启动各个协程
	go c.readPump()
	go c.writePump()
	go c.chatPump() //for test

	// 等待中断信号或连接断开
	select {
	case <-c.done:
		logx.Infof("[%s] 连接已断开", c.clientID)
	case <-c.interrupt:
		logx.Infof("[%s] 收到中断信号，正在关闭连接...", c.clientID)
	}

	c.Close()
}

// 读取消息协程
func (c *WSClient) readPump() {
	defer func() {
		c.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Infof("[%s] 读取消息错误: %v", c.clientID, err)
			}
			break
		}

		c.handleMessage(msg)
	}
}

// 写入协程 - 唯一的写入点
func (c *WSClient) writePump() {
	// 心跳定时器
	ticker := time.NewTicker(c.heartbeatInterval)
	defer func() {
		ticker.Stop()
		c.Close()
		logx.Infof("客户端 %s 写入协程退出", c.clientID)
	}()

	for {
		select {
		case message, ok := <-c.send:
			// 设置写入超时
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				// 发送通道已关闭，发送关闭消息
				logx.Infof("客户端 %s 发送通道关闭", c.clientID)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 写入消息
			if err := c.conn.WriteJSON(message); err != nil {
				logx.Infof("客户端 %s 写入消息失败: %v", c.clientID, err)
				return
			}

			logx.Infof("客户端 %s 发送消息: [%s] %s", c.clientID, message.Type, message.Content)

		case <-ticker.C:
			// 发送心跳 Ping
			heartbeat := Message{
				Type:      MsgTypeHeartbeat,
				Content:   "ping",
				Timestamp: time.Now(),
			}

			if !c.sendMessage(heartbeat) {
				logx.Infof("[%s] 发送心跳失败: %v", c.clientID)
				return
			}

			c.mutex.Lock()
			c.lastHeartbeat = time.Now()
			c.mutex.Unlock()

			logx.Infof("[%s] 发送心跳", c.clientID)

		case <-c.done:
			logx.Infof("客户端 %s 收到关闭信号", c.clientID)
			return
		}
	}
}

// 模拟聊天消息发送协程
func (c *WSClient) chatPump() {
	ticker := time.NewTicker(45 * time.Second) // 每45秒发送一条聊天消息
	defer ticker.Stop()

	messages := []string{
		"Hello from client!",
		"How is everyone doing?",
		"This is a test message",
		"WebSocket is working great!",
		"Hope everyone is having a good day",
	}

	for {
		select {
		case <-ticker.C:
			if !c.IsConnected() {
				return
			}

			// 随机选择一条消息
			content := messages[rand.Intn(len(messages))]

			chatMsg := Message{
				Type:      MsgTypeChat,
				Content:   fmt.Sprintf("[%s] %s", c.clientID, content),
				Timestamp: time.Now(),
			}

			if !c.sendMessage(chatMsg) {
				logx.Infof("[%s] 发送聊天消息失败: %v", c.clientID)
				return
			}

			logx.Infof("[%s] 发送聊天消息: %s", c.clientID, content)

		case <-c.done:
			return
		}
	}
}

// 处理接收到的消息
func (c *WSClient) handleMessage(msg Message) {
	switch msg.Type {
	case MsgTypeHeartbeatAck:
		logx.Infof("[%s] 收到心跳确认", c.clientID)

	case MsgTypeSystem:
		logx.Infof("[%s] 系统消息: %s", c.clientID, msg.Content)

	case MsgTypeChat:
		if msg.ClientID != c.clientID { // 不显示自己的消息
			logx.Infof("[%s] 聊天消息: %s", c.clientID, msg.Content)
		}

	default:
		logx.Infof("[%s] 未知消息类型 [%s]: %s", c.clientID, msg.Type, msg.Content)
	}
}

// 发送消息
func (c *WSClient) sendMessage(msg Message) bool {
	//c.mutex.RLock()
	//defer c.mutex.RUnlock()
	//c.mutex.Lock()
	//defer c.mutex.Unlock()
	//if !c.isConnected || c.conn == nil {
	//	return fmt.Errorf("连接已断开")
	//}
	//
	//c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	//return c.conn.WriteJSON(msg)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		logx.Infof("客户端 %s 已关闭，无法发送消息", c.clientID)
		return false
	}

	select {
	case c.send <- msg:
		return true
	case <-time.After(5 * time.Second): // 防止阻塞
		logx.Infof("客户端 %s 发送消息超时", c.clientID)
		return false
	}

}

// 检查连接状态
func (c *WSClient) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// 关闭连接
func (c *WSClient) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		c.isConnected = false

		// 发送关闭消息
		if c.conn != nil {
			c.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			c.conn.Close()
		}

		// 通知所有协程退出
		select {
		case <-c.done:
			// 已经关闭
		default:
			close(c.done)
		}

		logx.Infof("[%s] 连接已关闭", c.clientID)
	}
}

// 获取最后心跳时间
func (c *WSClient) GetLastHeartbeat() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lastHeartbeat
}
