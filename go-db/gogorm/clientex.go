package gogorm

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/gorm"
)

var _clientsForDbType = map[string]*GoGorm{}

func InitByDbTypeWithGormConfig(gormConfig gorm.Config, configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		key := conf.DBType + ":" + name
		if _clientsForDbType[key] != nil {
			return errors.New("gogorm clientForDbType already exists")
		}

		_clientsForDbType[key], err = New(&conf, gormConfig)
		if err != nil {
			return
		}
	}

	return
}

func InitByDbType(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		key := conf.DBType + ":" + name
		if _clientsForDbType[key] != nil {
			return errors.New("gogorm clientForDbType already exists")
		}

		_clientsForDbType[key], err = New(&conf)
		if err != nil {
			return
		}
	}

	return
}

func GetClientByDbType(dbType string, names ...string) *GoGorm {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	key := dbType + ":" + name
	if cli, ok := _clientsForDbType[key]; ok {
		return cli
	}
	return nil
}

func DelClientDbType(dbType string, names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			key := dbType + ":" + name
			delete(_clientsForDbType, key)
		}
	}
}

func DefaultForDbType(dbType string) *GoGorm {
	if cli, ok := _clientsForDbType[dbType+":default"]; ok {
		return cli
	}

	if l := len(_clientsForDbType); l == 1 {
		for _, cli := range _clientsForDbType {
			return cli
		}
	}

	golog.WithTag("gogorm").Error("no default gogorm clientForDbType")

	return nil
}
