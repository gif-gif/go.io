package goutils

import (
	"sync"
	"time"
)

//缺点：但是雪花算法强依赖机器时钟，如果机器上时钟回拨，会导致发号重复或者服务会处于不可用状态。如果恰巧回退前生成过一些ID，而时间回退后，生成的ID就有可能重复。官方对于此并没有给出解决方案，而是简单的抛错处理，这样会造成在时间被追回之前的这段时间服务不可用。很多其他类雪花算法也是在此思想上的设计然后改进规避它的缺陷，
//后面介绍的百度 UidGenerator 和 美团分布式ID生成系统 Leaf 中snowflake模式都是在 snowflake 的基础上演进出来的。#

// Snowflake，雪花算法是由Twitter开源的分布式ID生成算法，以划分命名空间的方式将 64-bit位分割成多个部分，
// 每个部分代表不同的含义。这种就是将64位划分为不同的段，每段代表不同的涵义，基本就是时间戳、机器ID和序列数。
// 为什么如此重要？因为它提供了一种ID生成及生成的思路，当然这种方案就是需要考虑时钟回拨的问题以及做一些 buffer的缓冲设计提高性能。
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
