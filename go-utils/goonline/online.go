package goonline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gif-gif/go.io/go-utils/gojson"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// GoOnline 通用在线管理器
type GoOnline struct {
	client     *clientv3.Client
	entityType string // 实体类型（users/servers/devices等）
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	Config     Config
	LogEnable  bool
}

// New 创建在线管理器
//
// entityType: 实体类型，如 "users", "servers", "devices" 等
//
// eg: /online/servers/serverId
func New(client *clientv3.Client, cfg *Config) (*GoOnline, error) {
	if cfg.OnlinePrefix == "" {
		cfg.OnlinePrefix = DefaultOnlinePrefix
	}
	// 设置默认租约过期时间
	if cfg.LeaseTTL <= 0 {
		cfg.LeaseTTL = DefaultLeaseTTL
	}
	// 设置默认超时时间
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultTimeout
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &GoOnline{
		client:     client,
		entityType: cfg.EntityType,
		ctx:        ctx,
		cancel:     cancel,
		Config:     *cfg,
		LogEnable:  cfg.LogEnable,
	}, nil
}

// getKey 获取实体的 etcd key
func (m *GoOnline) getKey(entityID string) string {
	return fmt.Sprintf("%s%s/%s", m.Config.OnlinePrefix, m.entityType, entityID)
}

// getPrefix 获取该类型实体的键前缀
func (m *GoOnline) getPrefix() string {
	return fmt.Sprintf("%s%s/", m.Config.OnlinePrefix, m.entityType)
}

// timeoutCtx 创建带超时的 context，避免 etcd 调用无限期阻塞
func (m *GoOnline) timeoutCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(m.ctx, time.Duration(m.Config.Timeout)*time.Second)
}

// SetOnline 设置实体上线或续租
// 如果实体已在线，则续租现有租约；如果不在线，则创建新租约
func (m *GoOnline) SetOnline(entityID string, leaseTTL int64, data any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := m.getKey(entityID)

	// 1. 检查实体是否已在线
	ctx, cancel := m.timeoutCtx()
	defer cancel()
	resp, err := m.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("查询实体状态失败: %v", err)
	}

	// 实体已在线，续租现有租约
	if len(resp.Kvs) > 0 {
		var existingData OnlineData
		err = gojson.Unmarshal(resp.Kvs[0].Value, &existingData)
		if err != nil {
			return fmt.Errorf("解析数据失败: %v", err)
		}

		// 续租现有租约
		leaseID := clientv3.LeaseID(existingData.LeaseID)
		kaCtx, kaCancel := m.timeoutCtx()
		defer kaCancel()
		_, err = m.client.KeepAliveOnce(kaCtx, leaseID)
		if err != nil {
			// 日志
			if m.LogEnable {
				logx.Infof("⚠ [%s]续租失败，租约可能已过期，将创建新租约: %v", entityID, err)
			}
			return m.createNewLease(entityID, leaseTTL, data)
		}

		// 更新数据
		existingData.Data = data
		existingData.UpdateAt = time.Now().Format(time.RFC3339)

		dataJSON, err := gojson.Marshal(existingData)
		if err != nil {
			return fmt.Errorf("序列化数据失败: %v", err)
		}

		// 更新 etcd
		putCtx, putCancel := m.timeoutCtx()
		defer putCancel()
		_, err = m.client.Put(putCtx, key, string(dataJSON), clientv3.WithLease(leaseID))
		if err != nil {
			return fmt.Errorf("更新数据失败: %v", err)
		}

		// 日志
		if m.LogEnable {
			logx.Infof("🔄 [%s] %s 续租成功，租约ID: %d", m.entityType, entityID, leaseID)
		}
		return nil
	}

	// 实体不在线，创建新租约
	return m.createNewLease(entityID, leaseTTL, data)
}

// createNewLease 创建新租约
func (m *GoOnline) createNewLease(entityID string, leaseTTL int64, data any) error {
	// 1. 创建租约
	if leaseTTL == 0 {
		leaseTTL = m.Config.LeaseTTL
	}
	grantCtx, grantCancel := m.timeoutCtx()
	defer grantCancel()
	lease, err := m.client.Grant(grantCtx, leaseTTL)
	if err != nil {
		return fmt.Errorf("创建租约失败: %v", err)
	}

	// 2. 构造数据
	now := time.Now().Format(time.RFC3339)
	onlineData := OnlineData{
		ID:       entityID,
		Type:     m.entityType,
		LeaseID:  int64(lease.ID),
		Data:     data,
		OnlineAt: now,
		UpdateAt: now,
	}

	dataJSON, err := gojson.Marshal(onlineData)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	// 3. 存储到 etcd
	key := m.getKey(entityID)
	putCtx, putCancel := m.timeoutCtx()
	defer putCancel()
	_, err = m.client.Put(putCtx, key, string(dataJSON), clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("注册失败: %v", err)
	}

	// 日志
	if m.LogEnable {
		logx.Infof("✓ [%s] %s 上线成功，租约ID: %d，过期时间: %d秒", m.entityType, entityID, lease.ID, leaseTTL)
	}
	return nil
}

// SetOffline 设置实体下线
func (m *GoOnline) SetOffline(entityID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := m.getKey(entityID)

	// 获取实体信息
	getCtx, getCancel := m.timeoutCtx()
	defer getCancel()
	resp, err := m.client.Get(getCtx, key)
	if err != nil {
		return fmt.Errorf("查询失败: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("[%s] %s 不在线", m.entityType, entityID)
	}

	// 解析数据获取租约ID
	var onlineData OnlineData
	err = gojson.Unmarshal(resp.Kvs[0].Value, &onlineData)
	if err != nil {
		return fmt.Errorf("解析数据失败: %v", err)
	}

	// 撤销租约
	revokeCtx, revokeCancel := m.timeoutCtx()
	defer revokeCancel()
	_, err = m.client.Revoke(revokeCtx, clientv3.LeaseID(onlineData.LeaseID))
	if err != nil {
		return fmt.Errorf("撤销租约失败: %v", err)
	}

	// 日志
	if m.LogEnable {
		logx.Infof("✓ [%s] %s 下线成功", m.entityType, entityID)
	}
	return nil
}

// GetOnlineList 获取所有在线实体列表
func (m *GoOnline) GetOnlineList() ([]OnlineData, error) {
	prefix := m.getPrefix()
	ctx, cancel := m.timeoutCtx()
	defer cancel()
	resp, err := m.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("获取在线列表失败: %v", err)
	}

	list := make([]OnlineData, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var data OnlineData
		err := gojson.Unmarshal(kv.Value, &data)
		if err != nil {
			// 日志
			if m.LogEnable {
				logx.Infof("解析数据失败: %v", err)
			}
			continue
		}
		list = append(list, data)
	}

	return list, nil
}

// Get 获取指定实体信息
func (m *GoOnline) Get(entityID string) (*OnlineData, error) {
	key := m.getKey(entityID)
	ctx, cancel := m.timeoutCtx()
	defer cancel()
	resp, err := m.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("[%s] %s 不在线", m.entityType, entityID)
	}

	var data OnlineData
	err = gojson.Unmarshal(resp.Kvs[0].Value, &data)
	if err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return &data, nil
}

// GetOnlineCount 获取在线数量
func (m *GoOnline) GetOnlineCount() (int, error) {
	prefix := m.getPrefix()
	ctx, cancel := m.timeoutCtx()
	defer cancel()
	resp, err := m.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, fmt.Errorf("获取在线数量失败: %v", err)
	}
	return int(resp.Count), nil
}

// IsOnline 检查是否在线
func (m *GoOnline) IsOnline(entityID string) (bool, error) {
	key := m.getKey(entityID)
	ctx, cancel := m.timeoutCtx()
	defer cancel()
	resp, err := m.client.Get(ctx, key, clientv3.WithCountOnly())
	if err != nil {
		return false, err
	}
	return resp.Count > 0, nil
}

// Close 关闭管理器, 谨慎调用, 会导致所有实体下线，且无法重新上线，建议仅在测试环境使用
func (m *GoOnline) Close() error {
	m.cancel()
	return m.client.Close()
}

// OnlineWatcher 通用在线监听器
type OnlineWatcher struct {
	client     *clientv3.Client
	entityType string // 实体类型
	ctx        context.Context
	cancel     context.CancelFunc
	Config     Config
}

// NewOnlineWatcher 创建在线监听器
func NewOnlineWatcher(client *clientv3.Client, config *Config) (*OnlineWatcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &OnlineWatcher{
		client:     client,
		entityType: config.EntityType,
		ctx:        ctx,
		cancel:     cancel,
		Config:     *config,
	}, nil
}

// getPrefix 获取键前缀
func (w *OnlineWatcher) getPrefix() string {
	return fmt.Sprintf("%s%s/", w.Config.OnlinePrefix, w.entityType)
}

// timeoutCtx 创建带超时的 context
func (w *OnlineWatcher) timeoutCtx() (context.Context, context.CancelFunc) {
	timeout := w.Config.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return context.WithTimeout(w.ctx, time.Duration(timeout)*time.Second)
}

// GetOnlineList 获取当前在线列表
func (w *OnlineWatcher) GetOnlineList() ([]OnlineData, error) {
	prefix := w.getPrefix()
	ctx, cancel := w.timeoutCtx()
	defer cancel()
	resp, err := w.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("获取在线列表失败: %v", err)
	}

	list := make([]OnlineData, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var data OnlineData
		err := gojson.Unmarshal(kv.Value, &data)
		if err != nil {
			// 日志
			if w.Config.LogEnable {
				logx.Errorf("解析数据失败: %v", err)
			}
			continue
		}
		list = append(list, data)
	}

	return list, nil
}

// Watch 监听在线状态变化 for test
func (w *OnlineWatcher) Watch() {
	// 日志
	if w.Config.LogEnable {
		logx.Infof("📡 开始监听 [%s] 在线状态变化...", w.entityType)
	}

	// 先获取当前在线列表
	list, err := w.GetOnlineList()
	if err != nil {
		// 日志
		if w.Config.LogEnable {
			logx.Errorf("获取初始在线列表失败: %v", err)
		}
		return
	}

	// 日志
	if w.Config.LogEnable {
		logx.Infof("📊 当前 [%s] 在线数: %d", w.entityType, len(list))
		for _, item := range list {
			logx.Infof("  👤 %s: %v", item.ID, item.Data)
		}
	}

	// 监听后续变化
	prefix := w.getPrefix()
	watchChan := w.client.Watch(w.ctx, prefix, clientv3.WithPrefix())

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			var data OnlineData
			if event.Type == clientv3.EventTypePut {
				gojson.Unmarshal(event.Kv.Value, &data)
				// 日志
				if w.Config.LogEnable {
					logx.Infof("🟢 [%s] [上线/更新] %s: %v", w.entityType, data.ID, data.Data)
				}
			} else if event.Type == clientv3.EventTypeDelete {
				// 从 key 中提取 ID
				key := string(event.Kv.Key)
				prefix := w.getPrefix()
				entityID := key[len(prefix):]
				// 日志
				if w.Config.LogEnable {
					logx.Infof("🔴 [%s] [下线] %s", w.entityType, entityID)
				}
			}
		}
	}
}

// Close 关闭监听器
func (w *OnlineWatcher) Close() error {
	w.cancel()
	return w.client.Close()
}
