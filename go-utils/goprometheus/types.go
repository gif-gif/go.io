package goprometheus

import "github.com/prometheus/common/model"

type MetricName = string
type MetricLabel = model.LabelName

type MetricQuery = struct {
	ProductCode string
	Group       string
	InstanceIds []int64
	TimeRange   string
}

type Bandwidth struct {
	In  int64 `json:"in,optional"`
	Out int64 `json:"out,optional"`
}

type SysUsage struct {
	CpuUsage    float64 `json:"cpuUsage,optional"`    // CPU使用率
	DiskTotal   int64   `json:"diskTotal,optional"`   // 磁盘总大小
	DiskUsage   float64 `json:"diskUsage,optional"`   // 磁盘使用率
	MemoryTotal int64   `json:"memoryTotal,optional"` // 内存总大小
	MemoryUsage float64 `json:"memoryUsage,optional"` // 内存使用率
}

type Traffic struct {
	In    int64 `json:"in,optional"`    // 入站流量
	Out   int64 `json:"out,optional"`   // 出站流量
	Total int64 `json:"total,optional"` // 总流量
}

// 会员等级用户数
type MemberLevelUserCount struct {
	Level int64 `json:"level,optional"` // 会员等级
	Count int64 `json:"count,optional"` // 用户数
} // 会员等级用户数

type UserLevelUserCount struct {
	Level int64 `json:"level,optional"` // 会员等级
	Count int64 `json:"count,optional"` // 用户数
}

type IpUserCount struct {
	Ip    string `json:"ip,optional"`    // IP地址
	Count int64  `json:"count,optional"` // 用户数
}

type IpBandwidth struct {
	Ip      string               `json:"ip,optional"`                // IP地址
	Details []*ProtocolBandwidth `json:"details,optional,omitempty"` // 带宽详情
}

type IpConnCount struct {
	Ip      string               `json:"ip,optional"`                // IP地址
	Details []*ProtocolConnCount `json:"details,optional,omitempty"` // 连接数详情
}

type ProtocolBandwidth struct {
	Protocol string `json:"protocol,optional"` // 协议
	Port     int64  `json:"port,optional"`     // 端口
	In       int64  `json:"in,optional"`       // 入站带宽
	Out      int64  `json:"out,optional"`      // 出站带宽
}

type ProtocolConnCount struct {
	Protocol string `json:"protocol,optional"` // 协议
	Port     int64  `json:"port,optional"`     // 端口
	In       int64  `json:"in,optional"`       // 入站连接数
	Out      int64  `json:"out,optional"`      // 出站连接数
}
