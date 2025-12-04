package govault

import (
	"fmt"
	"net/url"

	"github.com/gogf/gf/util/gconv"
)

type UserNameAndPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 默认 charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
type MysqlDataSource struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       int64  `json:"port"`
	Database   string `json:"database"`
	Charset    string `json:"charset,optional" default:"utf8mb4"`
	ParseTime  string `json:"parseTime,optional" default:"true"`
	Loc        string `json:"loc,optional" default:"Asia%2FShanghai"`
	DataSource string `json:"dataSource,optional"`
}

type MysqlDataSourceSimple struct {
	DataSource string `json:"dataSource"`
}

func ParseData(data map[string]interface{}) MysqlDataSource {
	port := gconv.Int64(data["port"])
	charset := gconv.String(data["charset"])
	loc := gconv.String(data["loc"])
	parseTime := gconv.String(data["parseTime"])
	if charset == "" {
		charset = "utf8mb4"
	}
	if parseTime == "" {
		parseTime = "true"
	}
	if loc == "" {
		loc = url.PathEscape("Asia/Shanghai")
	}
	var ds = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s", data["username"].(string), data["password"].(string), data["host"].(string), port, data["database"].(string), charset, parseTime, loc)
	return MysqlDataSource{
		DataSource: ds,
		Username:   data["username"].(string),
		Password:   data["password"].(string),
		Host:       data["host"].(string),
		Port:       port,
		Database:   data["database"].(string),
		Charset:    data["charset"].(string),
		ParseTime:  data["parseTime"].(string),
		Loc:        data["loc"].(string),
	}
}

func ParseMap(data MysqlDataSource) map[string]interface{} {
	return map[string]interface{}{
		"username":  data.Username,
		"password":  data.Password,
		"host":      data.Host,
		"port":      data.Port,
		"database":  data.Database,
		"charset":   data.Charset,
		"parseTime": data.ParseTime,
		"loc":       data.Loc,
	}
}
