package goonline

const (
	DefaultLeaseTTL     = 600        // 租约过期时间：10分钟
	DefaultOnlinePrefix = "/online/" // 在线数据键前缀
)

// Config 配置项
type Config struct {
	//Etcd         goetcd.Config
	Name         string `yaml:"name" json:"name,optional"`
	EntityType   string `yaml:"entityType" json:"entityType,optional"`     // 实体类型（users/servers/devices等）
	LeaseTTL     int64  `yaml:"leaseTTL" json:"leaseTTL,optional"`         // 租约过期时间（秒）,默认10分钟
	OnlinePrefix string `yaml:"onlinePrefix" json:"onlinePrefix,optional"` // 在线数据键前缀, 默认"/online/"
	LogEnable    bool   `yaml:"logEnable" json:"logEnable,optional"`       // 是否开启日志
}

// OnlineData 存储在 etcd 中的通用在线数据结构
type OnlineData struct {
	ID       string `json:"id"`        // 实体ID
	Type     string `json:"type"`      // 实体类型（user/server/device等）
	LeaseID  int64  `json:"lease_id"`  // 租约ID
	Data     any    `json:"data"`      // 自定义数据
	OnlineAt string `json:"online_at"` // 上线时间
	UpdateAt string `json:"update_at"` // 最后更新时间
}
