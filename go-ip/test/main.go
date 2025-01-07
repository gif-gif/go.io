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
		Mmdb:          "/Users/Jerry/Documents/my/dockers/projects/golang/ip_service/ip/data/GeoLite2-Country.mmdb",
		Ip2locationDB: "/Users/Jerry/Documents/my/dockers/projects/golang/ip_service/ip/data/IP-COUNTRY-REGION-CITY-ISP.BIN",
		IpServiceUrl:  "",
	}

	ip := "154.18.180.43"
	err := goip.Init(config)
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	ipZhInfo, err := goip.Default().QueryDbReaderCountryForZhName(ip)
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipZhInfo)

	ipInfo, err := goip.Default().QueryLocationInfoByIp(ip)
	if err != nil {
		golog.WithTag("goip").Error(err.Error())
		return
	}

	golog.WithTag("goip").Info(ipInfo)

	if config.IpServiceUrl != "" {
		ipinfo, err := goip.Default().GetIpLocation(context.Background(), ip)
		if err != nil {
			golog.WithTag("goip").Error(err.Error())
			return
		}

		golog.WithTag("goip").Info(ipinfo)
	}

	time.Sleep(time.Second * 5)

}
