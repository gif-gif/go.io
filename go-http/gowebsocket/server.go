package gowebsocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServerConfig struct {
	Name                            string `yaml:"Name" json:"name,optional"`
	Port                            int    `yaml:"Port" json:"port,optional"`
	HandshakeTimeout                time.Duration
	ReadBufferSize, WriteBufferSize int
	WriteBufferPool                 websocket.BufferPool
	Subprotocols                    []string
	Error                           func(w http.ResponseWriter, r *http.Request, status int, reason error)
	CheckOrigin                     func(r *http.Request) bool
	EnableCompression               bool

	HeartbeatInterval int64 `yaml:"HeartbeatInterval" json:"heartbeatInterval,optional"` //秒
	HeartbeatTimeout  int64 `yaml:"HeartbeatTimeout" json:"heartbeatTimeout,optional"`   //秒
	ShutdownTimeout   int64 `yaml:"ShutdownTimeout" json:"shutdownTimeout,optional"`     // 秒
}

type Client struct {
	ID            string
	Conn          *websocket.Conn
	Send          chan Message
	Hub           *Hub
	LastHeartbeat time.Time
	IsAlive       bool
	mutex         sync.RWMutex
}

type Server struct {
	hub        *Hub
	httpServer *http.Server
	wg         sync.WaitGroup
	shutdown   chan os.Signal
}

type Hub struct {
	clients        map[string]*Client
	register       chan *Client
	unregister     chan *Client
	broadcast      chan Message
	shutdown       chan struct{}
	done           chan struct{}
	mutex          sync.RWMutex
	isShuttingDown bool
	upgrader       *websocket.Upgrader
	Config         ServerConfig
}

// 创建新的 Hub
func NewHub(config ServerConfig) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
		shutdown:   make(chan struct{}),
		done:       make(chan struct{}),
		Config:     config,
	}
}

func create(config ServerConfig) (ws *Server) {
	c := &config
	if c.HeartbeatInterval <= 0 {
		c.HeartbeatInterval = 30
	}
	if c.HeartbeatTimeout <= 0 {
		c.HeartbeatTimeout = 60
	}
	if c.ShutdownTimeout <= 0 {
		c.ShutdownTimeout = 30
	}
	h := NewHub(config)
	ws = &Server{
		hub: h,
	}
	h.upgrader = &websocket.Upgrader{
		ReadBufferSize:    config.ReadBufferSize,
		WriteBufferSize:   config.WriteBufferSize,
		Subprotocols:      config.Subprotocols,
		Error:             config.Error,
		EnableCompression: config.EnableCompression,
		CheckOrigin:       config.CheckOrigin,
	}
	return ws
}

// 创建新的服务器
func NewServer(config ServerConfig) *Server {
	server := create(config)
	hub := server.hub

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.wsHandler(w, r)
	})
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		statsHandler(hub, w, r)
	})

	server.httpServer = &http.Server{
		Addr:    gconv.String(":" + gconv.String(config.Port)),
		Handler: mux,
	}
	server.shutdown = make(chan os.Signal, 1)
	// 监听系统信号
	signal.Notify(server.shutdown, os.Interrupt, syscall.SIGTERM)
	return server
}

// 启动服务器
func (s *Server) Start() error {
	logx.Infof("WebSocket 服务器启动在 %s", s.httpServer.Addr)
	logx.Infof("心跳检查间隔: %d秒", s.hub.Config.HeartbeatInterval)
	logx.Infof("心跳超时时间: %d秒", s.hub.Config.HeartbeatTimeout)
	logx.Infof("关闭超时时间: %d秒", s.hub.Config.ShutdownTimeout)
	logx.Infof("按 Ctrl+C 优雅关闭服务器")

	// 启动 Hub
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.hub.RunMainLoop()
	}()

	// 启动统计协程
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.statsReporter()
	}()

	// 启动 HTTP 服务器
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Infof("HTTP 服务器错误: %v", err)
		}
	}()

	// 等待关闭信号
	<-s.shutdown
	logx.Infof("收到关闭信号，开始优雅关闭...")

	return s.GracefulShutdown(time.Duration(s.hub.Config.ShutdownTimeout) * time.Second)
}

// 健康检查处理器
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// 统计信息处理器
func statsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := hub.GetStats()
	json.NewEncoder(w).Encode(stats)
}

// 统计报告器
func (s *Server) statsReporter() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !s.hub.isShuttingDown {
				logx.Infof("当前在线客户端数量: %d", s.hub.GetOnlineCount())
			}
		case <-s.hub.done:
			return
		}
	}
}

// Hub 运行主循环
func (h *Hub) RunMainLoop() {
	defer close(h.done)

	// 启动心跳检查协程
	heartbeatTicker := time.NewTicker(time.Duration(h.Config.HeartbeatInterval) * time.Second)
	defer heartbeatTicker.Stop()

	logx.Infof("Hub 开始运行")

	for {
		select {
		case client := <-h.register:
			if h.isShuttingDown {
				// 服务器正在关闭，拒绝新连接
				logx.Infof("服务器关闭中，拒绝新连接: %s", client.ID)
				client.Conn.Close()
				continue
			}
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.BroadcastMessage(message)

		case <-heartbeatTicker.C:
			if !h.isShuttingDown {
				h.checkHeartbeats()
			}

		case <-h.shutdown:
			logx.Infof("Hub 收到关闭信号")
			h.mutex.Lock()
			h.isShuttingDown = true
			h.mutex.Unlock()

			// 关闭所有客户端连接
			h.closeAllClients()
			return
		}
	}
}

// 关闭所有客户端
func (h *Hub) closeAllClients() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	logx.Infof("开始关闭 %d 个客户端连接", len(h.clients))

	for clientID, client := range h.clients {
		logx.Infof("关闭客户端: %s", clientID)

		// 关闭发送通道
		close(client.Send)

		// 发送 WebSocket 关闭帧
		client.Conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "服务器关闭"),
		)

		// 关闭连接
		client.Conn.Close()

		delete(h.clients, clientID)
	}

	logx.Infof("所有客户端连接已关闭")
}

// 关闭 Hub
func (s *Server) Shutdown() {
	select {
	case <-s.hub.shutdown:
		// 已经在关闭
	default:
		close(s.hub.shutdown)
	}

	// 等待 Hub 协程完成
	<-s.hub.done
	logx.Infof("Hub 已关闭")
}

// 优雅关闭服务器
func (s *Server) GracefulShutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logx.Infof("=== 开始优雅关闭流程 ===")

	// 1. 停止接受新连接
	logx.Infof("1. 停止接受新的 HTTP 连接...")
	shutdownErr := s.httpServer.Shutdown(ctx)
	if shutdownErr != nil {
		logx.Infof("HTTP 服务器关闭错误: %v", shutdownErr)
	}

	// 2. 通知所有客户端服务器即将关闭
	logx.Infof("2. 通知所有客户端服务器关闭 (当前连接数: %d)...", s.hub.GetOnlineCount())
	s.hub.notifyShutdown()

	// 3. 给客户端一些时间处理关闭通知
	logx.Infof("3. 等待客户端处理关闭通知...")
	time.Sleep(3 * time.Second)

	// 4. 开始关闭 Hub
	logx.Infof("4. 关闭 Hub 和所有 WebSocket 连接...")
	s.Shutdown()

	// 5. 等待所有协程完成
	logx.Infof("5. 等待所有协程完成...")
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logx.Infof("6. 所有协程已完成")
	case <-ctx.Done():
		logx.Infof("6. 超时，强制关闭")
		return fmt.Errorf("优雅关闭超时: %v", ctx.Err())
	}

	logx.Infof("=== 服务器已优雅关闭 ===")
	return shutdownErr
}

// 通知关闭
func (h *Hub) notifyShutdown() {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	shutdownMsg := Message{
		Type:      MsgTypeShutdown,
		Content:   "服务器即将关闭，请准备断开连接",
		Timestamp: time.Now(),
	}

	logx.Infof("向 %d 个客户端发送关闭通知", len(h.clients))

	for _, client := range h.clients {
		select {
		case client.Send <- shutdownMsg:
		default:
			// 客户端发送队列满，直接关闭
			logx.Infof("客户端 %s 发送队列满，直接关闭", client.ID)
		}
	}
}

// WebSocket 处理函数
func (h *Hub) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Infof("WebSocket 升级失败: %v", err)
		return
	}

	// 生成客户端ID
	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())

	client := &Client{
		ID:            clientID,
		Conn:          conn,
		Send:          make(chan Message, 256),
		Hub:           h,
		LastHeartbeat: time.Now(), // 初始化心跳时间
		IsAlive:       true,
	}

	// 注册客户端
	client.Hub.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	h.clients[client.ID] = client
	h.mutex.Unlock()

	logx.Infof("客户端 %s 已连接，当前在线: %d", client.ID, len(h.clients))

	// 发送欢迎消息
	welcomeMsg := Message{
		Type:      MsgTypeSystem,
		Content:   fmt.Sprintf("欢迎 %s 加入服务器", client.ID),
		Timestamp: time.Now(),
	}

	select {
	case client.Send <- welcomeMsg:
	default:
		close(client.Send)
		h.mutex.Lock()
		delete(h.clients, client.ID)
		h.mutex.Unlock()
	}
}

// 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	if _, exists := h.clients[client.ID]; exists {
		delete(h.clients, client.ID)
		close(client.Send)
		logx.Infof("客户端 %s 已断开，当前在线: %d", client.ID, len(h.clients))
	}
	h.mutex.Unlock()
}

// 广播消息
func (h *Hub) BroadcastMessage(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if h.isShuttingDown {
		return // 关闭期间不广播消息
	}

	logx.Infof("广播消息: [%s] %s", message.Type, message.Content)

	for clientID, client := range h.clients {
		select {
		case client.Send <- message:
		default:
			// 发送失败，移除客户端
			logx.Infof("向客户端 %s 发送消息失败，移除连接", clientID)
			close(client.Send)
			delete(h.clients, clientID)
		}
	}
}

// 心跳检查器
func (h *Hub) checkHeartbeats() {
	h.mutex.RLock()
	var offlineClients []*Client

	for _, client := range h.clients {
		client.mutex.RLock()
		timeSinceLastHeartbeat := time.Since(client.LastHeartbeat)
		client.mutex.RUnlock()

		// 如果超过60秒没有收到心跳，认为客户端离线
		if timeSinceLastHeartbeat > time.Duration(h.Config.HeartbeatTimeout)*time.Second {
			logx.Infof("客户端 %s 心跳超时 (%v)，标记为离线",
				client.ID, timeSinceLastHeartbeat)
			offlineClients = append(offlineClients, client)
		}
	}
	h.mutex.RUnlock()

	// 移除离线客户端
	for _, client := range offlineClients {
		client.mutex.Lock()
		client.IsAlive = false
		client.mutex.Unlock()

		h.unregister <- client
		client.Conn.Close()
	}
}

// 获取在线客户端数量
func (h *Hub) GetOnlineCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// 获取客户端统计信息
func (h *Hub) GetStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["online_count"] = len(h.clients)
	stats["is_shutting_down"] = h.isShuttingDown

	clientList := make([]string, 0, len(h.clients))
	for clientID := range h.clients {
		clientList = append(clientList, clientID)
	}
	stats["clients"] = clientList

	return stats
}

// client
// 分析关闭错误类型
func (c *Client) analyzeCloseError(err error) {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		logx.Infof("✅ 客户端 %s 正常关闭连接", c.ID)
	} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
		logx.Infof("👋 客户端 %s 页面跳转或刷新", c.ID)
	} else if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
		logx.Infof("⚡ 客户端 %s 异常关闭连接（网络中断或进程终止）", c.ID)
	} else if websocket.IsCloseError(err, websocket.CloseUnsupportedData) {
		logx.Infof("❌ 客户端 %s 不支持数据类型", c.ID)
	} else if websocket.IsUnexpectedCloseError(err,
		websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
		logx.Infof("🔥 客户端 %s 意外关闭: %v", c.ID, err)
	} else {
		logx.Infof("📋 客户端 %s 连接错误: %v", c.ID, err)
	}
}

// 客户端读取消息
// 客户端读取消息
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// 设置 Close 处理器（客户端主动关闭）
	c.Conn.SetCloseHandler(func(code int, text string) error {
		logx.Infof("客户端 %s 主动关闭连接 (代码: %d, 原因: %s)", c.ID, code, text)
		return nil
	})

	// 设置读取超时
	c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		return nil
	})

	for {
		var msg Message
		//err := c.Conn.ReadJSON(&msg)
		//if err != nil {
		//	if websocket.IsUnexpectedCloseError(err,
		//		websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		//		logx.Infof("WebSocket 读取错误 [%s]: %v", c.ID, err)
		//	}
		//	break
		//}

		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			// 详细的错误分析
			c.analyzeCloseError(err)
			break
		}

		msg.ClientID = c.ID
		msg.Timestamp = time.Now()

		// 处理消息
		c.handleMessage(msg)
	}
}

// 客户端写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				logx.Infof("写入消息失败 [%s]: %v", c.ID, err)
				return
			}

			// 特殊处理关闭消息
			if message.Type == MsgTypeShutdown {
				logx.Infof("已向客户端 %s 发送关闭通知", c.ID)
			}

		case <-ticker.C: //TODO: 服务器心跳
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 处理接收到的消息
func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case MsgTypeHeartbeat:
		// 收到心跳，更新最后心跳时间
		c.mutex.Lock()
		c.LastHeartbeat = time.Now()
		c.IsAlive = true
		c.mutex.Unlock()

		// 回复心跳确认
		ack := Message{
			ClientID:  c.ID,
			Type:      MsgTypeHeartbeatAck,
			Content:   msg.Content,
			Timestamp: time.Now(),
		}

		select {
		case c.Send <- ack:
		default:
			logx.Infof("发送心跳确认失败: 客户端 %s", c.ID)
		}

	case MsgTypeChat:
		logx.Infof("收到聊天消息 [%s]: %s", c.ID, msg.Content)
		// 广播聊天消息
		c.Hub.broadcast <- msg

	default:
		logx.Infof("收到未知类型消息 [%s]: %s", msg.Type, msg.Content)
	}
}
