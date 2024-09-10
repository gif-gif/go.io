package main

import (
	"context"
	goip "github.com/gif-gif/go.io/go-ip"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	config := goip.Config{
		//Mmdb:          "GeoLite2-Country.mmdb",
		//Ip2locationDB: "IP-COUNTRY-REGION-CITY-ISP.BIN",
		IpServiceUrl: "http://172.99.189.235:20030/ip/country/v2",
	}
	err := goip.Init(config)
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	ipinfo, err := goip.Default().GetIpLocation(context.Background(), "172.99.189.235")
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipinfo)

	ipZhInfo, err := goip.Default().QueryDbReaderCountryForZhName("172.99.189.235")
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipZhInfo)

	ipInfo, err := goip.Default().QueryLocationInfoByIp("172.99.189.235")
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipInfo)

	time.Sleep(time.Second * 5)

}
