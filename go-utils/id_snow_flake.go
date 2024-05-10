package goutils

import (
	"sync"
	"time"
)

// 雪花算法
type SnowFlakeId struct {
	dataCenterId int // 机房ID
	machineId    int // 机器ID

	lastTime int64 // 最后时间
	sn       int   // 序号

	mu sync.Mutex
}

func (sf *SnowFlakeId) GenId() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	ts := time.Now().UnixNano() / 1e6

	if sf.lastTime == ts {
		// 2的12次方 -1 = 4095，每毫秒可产生4095个ID
		if sf.sn > 4095 {
			time.Sleep(time.Millisecond)
			ts = time.Now().UnixNano() / 1e6
			sf.sn = 0
		}
	} else {
		sf.sn = 0
	}

	sf.sn += 1
	sf.lastTime = ts

	// 时间戳，向左移动22位
	ts = ts << 22

	// 机房ID，向左移动17位
	dataCenterId := sf.dataCenterId << 17

	// 机器ID，向左移动12位
	machineId := sf.machineId << 12

	return ts | int64(dataCenterId) | int64(machineId) | int64(sf.sn)
}
