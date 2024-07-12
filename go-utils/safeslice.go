package goutils

import (
	golock "github.com/gif-gif/go.io/go-lock"
	"sync"
)

// 读写锁 + 顺序获取（循环）
//
// 如： s := []string{"s1", "s2", "s3"}， 无论并发怎么读取，顺序为：s1,s2,s3,s1,s2,s3,s1,s2,s...
type SafeSlice[T comparable] struct {
	lock       *golock.GoLock
	currentKey T
	data       []T
}

func NewSafeSlice[T comparable]() *SafeSlice[T] {
	return &SafeSlice[T]{
		lock: &golock.GoLock{
			MuteRW: *new(sync.RWMutex),
		},
	}
}

func (m *SafeSlice[T]) Set(str T) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data = append(m.data, str)
	m.data = m.removeDuplicates(m.data)
}

func (m *SafeSlice[T]) Sets(data []T) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.data == nil {
		m.data = data
	} else {
		m.data = append(m.data, data...)
		m.data = m.removeDuplicates(m.data)
	}
}

func (m *SafeSlice[T]) removeDuplicates(s []T) []T {
	seen := make(map[interface{}]struct{})
	result := make([]T, 0, len(s))
	for _, value := range s {
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

func (m *SafeSlice[T]) Get() T {
	m.lock.Lock()
	defer m.lock.Unlock()
	count := len(m.data) - 1
	key := m.data[0]
	for k, v := range m.data {
		if v == m.currentKey {
			if k < count {
				key = m.data[k+1]
				break
			}
		}
	}
	m.currentKey = key
	return key
}
