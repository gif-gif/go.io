package goprometheus

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

type MemberLevelUserCount struct {
	Level int64 `json:"level,optional"` // 会员等级
	Count int64 `json:"count,optional"` // 用户数
}
