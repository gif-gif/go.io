package main

import (
	goip "github.com/gif-gif/go.io/go-ip"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	config := goip.Config{
		Mmdb:          "GeoLite2-Country.mmdb",
		Ip2locationDB: "IP-COUNTRY-REGION-CITY-ISP.BIN",
	}
	g := &goip.GoIp{}
	err := g.Init(config)
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	ipZhInfo, err := g.QueryDbReaderCountryForZhName("172.99.189.235")
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipZhInfo)

	ipInfo, err := g.QueryLocationInfoByIp("172.99.189.235")
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipInfo)

	time.Sleep(time.Second * 5)

}
