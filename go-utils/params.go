package goutils

import (
	"encoding/json"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Byte []byte

func (b Byte) Params() (p *Params, err error) {
	p = NewParams()
	if err = json.Unmarshal(b, &p.data); err != nil {
		goo_log.WithField("params", string(b)).Error(err)
		return
	}
	return
}

type Params struct {
	data interface{}
	mu   sync.RWMutex
}

func NewParams() *Params {
	return &Params{data: map[string]interface{}{}}
}

func (p *Params) Set(key string, val interface{}) *Params {
	p.mu.Lock()
	if p.data == nil {
		p.data = map[string]interface{}{}
	}
	if v, ok := (p.data).(map[string]interface{}); ok {
		v[key] = val
	}
	p.mu.Unlock()
	return p
}

func (p *Params) Get(key string) *Params {
	keys := strings.Split(key, ".")
	for _, k := range keys {
		if data, ok := (p.data).(map[string]interface{}); ok {
			if v, ok := data[k]; ok {
				p.data = v
				continue
			}
			p.data = nil
			break
		}
		p.data = nil
		break
	}
	return p
}

func (p *Params) String() string {
	switch reflect.ValueOf(p.data).Kind() {
	case reflect.Float64:
		return strconv.FormatFloat(p.Float64(), 'f', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(p.Float64(), 'f', -1, 64)
	case reflect.Int:
		return strconv.FormatInt(p.Int64(), 10)
	case reflect.Int32:
		return strconv.FormatInt(p.Int64(), 10)
	case reflect.Int64:
		return strconv.FormatInt(p.Int64(), 10)
	case reflect.Bool:
		return strconv.FormatBool(p.Bool())
	case reflect.String:
		return (p.data).(string)
	}
	return string(p.JSON())
}

func (p *Params) Int64() int64 {
	switch reflect.ValueOf(p.data).Kind() {
	case reflect.Float64:
		return int64((p.data).(float64))
	case reflect.Float32:
		return int64((p.data).(float32))
	case reflect.Int:
		return int64((p.data).(int))
	case reflect.Int32:
		return int64((p.data).(int32))
	case reflect.Int64:
		return (p.data).(int64)
	case reflect.Bool:
		if (p.data).(bool) {
			return 1
		}
	case reflect.String:
		v, _ := strconv.ParseInt((p.data).(string), 10, 64)
		return v
	}
	return 0
}

func (p *Params) Int32() int32 {
	return int32(p.Int64())
}

func (p *Params) Int() int {
	return int(p.Int64())
}

func (p *Params) Float64() float64 {
	if v, ok := (p.data).(float64); ok {
		return v
	}
	return 0
}

func (p *Params) Float32() float32 {
	if v, ok := (p.data).(float32); ok {
		return v
	}
	return 0
}

func (p *Params) Bool() bool {
	if v, ok := (p.data).(bool); ok {
		return v
	}
	return false
}

func (p *Params) Array() (ps []Params) {
	ps = []Params{}
	if arr, ok := (p.data).([]interface{}); ok {
		for _, data := range arr {
			ps = append(ps, Params{data: data})
		}
	}
	return
}

func (p *Params) Map() (rst map[string]Params) {
	rst = map[string]Params{}
	if m, ok := (p.data).(map[string]interface{}); ok {
		for k, data := range m {
			rst[k] = Params{data: data}
		}
	}
	return
}

func (p *Params) Data() interface{} {
	return p.data
}

func (p *Params) MapData() map[string]interface{} {
	if data, ok := (p.data).(map[string]interface{}); ok {
		return data
	}
	return map[string]interface{}{}
}

func (p *Params) ArrayData() []interface{} {
	if data, ok := (p.data).([]interface{}); ok {
		return data
	}
	return []interface{}{}
}

func (p *Params) JSON() []byte {
	buf, _ := json.Marshal(p.data)
	return buf
}
