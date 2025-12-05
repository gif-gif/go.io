package gosnow

import (
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type iGenId interface {
	GenId() int64
}

var (
	__genId    iGenId
	__genIdOne sync.Once
)

func GenIdInit(adapter iGenId) {
	__genId = adapter
}

func GenId() int64 {
	__genIdOne.Do(func() {
		if __genId == nil {
			__genId = &SnowFlakeId{WorkerId: 1}
		}
	})
	return __genId.GenId()
}

func GenIdStr() string {
	return strconv.FormatInt(GenId(), 10)
}

func UUID() string {
	return uuid.New().String()
}
