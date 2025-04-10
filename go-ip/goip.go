package goip

import (
	"context"
	"fmt"
	gocache "github.com/gif-gif/go.io/go-cache"
	gohttp "github.com/gif-gif/go.io/go-http"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/oschwald/geoip2-golang"
	"golang.org/x/sync/singleflight"
	"net"
	"strings"
	"time"
)

type GoIp struct {
	dbReader          *geoip2.Reader
	dbLocation        *ip2location.DB
	IpServiceUrl      string
	singlefightGForIp singleflight.Group
}

func New(config Config) (*GoIp, error) {
	g := &GoIp{
		IpServiceUrl:      config.IpServiceUrl,
		singlefightGForIp: singleflight.Group{},
	}
	if config.Mmdb != "" {
		db, err := geoip2.Open(config.Mmdb)
		if err != nil {
			return nil, err
		}
		g.dbReader = db
	}

	if config.Ip2locationDB != "" {
		dbLocation, err := ip2location.OpenDB(config.Ip2locationDB)
		if err != nil {
			return nil, err
		}
		g.dbLocation = dbLocation
	}

	return g, nil
}

// 查询 IP相关信息,返回IP所属国家中文名称
func (g *GoIp) QueryDbReaderCountryForZhName(ipStr string) (*IpCountry, error) {
	if g.dbReader == nil {
		return nil, fmt.Errorf("dbReader is nil")
	}
	ip := net.ParseIP(ipStr)
	recordCountry, err := g.dbReader.Country(ip)
	if err != nil {
		return nil, err
	}
	country := &IpCountry{
		IsoCode:       recordCountry.Country.IsoCode,
		Name:          recordCountry.Country.Names["zh-CN"],
		Continent:     recordCountry.Continent.Names["zh-CN"],
		ContinentCode: recordCountry.Continent.Code,
	}
	return country, nil
}

// IP 相关信息
func (g *GoIp) QueryLocationInfoByIp(ipStr string) (*IP2Locationrecord, error) {
	if g.dbLocation == nil {
		return nil, fmt.Errorf("dbLocation is nil")
	}
	results, err := g.dbLocation.Get_all(ipStr)
	if err != nil {
		return nil, err
	}

	rsp := &IP2Locationrecord{
		Country_short:      results.Country_short,
		Country_long:       results.Country_long,
		Region:             results.Region,
		City:               results.City,
		Isp:                results.Isp,
		Latitude:           results.Latitude,
		Longitude:          results.Longitude,
		Domain:             results.Domain,
		Zipcode:            results.Zipcode,
		Timezone:           results.Timezone,
		Netspeed:           results.Netspeed,
		Iddcode:            results.Iddcode,
		Areacode:           results.Areacode,
		Weatherstationcode: results.Weatherstationcode,
		Weatherstationname: results.Weatherstationname,
		Mcc:                results.Mcc,
		Mnc:                results.Mnc,
		Mobilebrand:        results.Mobilebrand,
		Elevation:          results.Elevation,
		Usagetype:          results.Usagetype,
		Addresstype:        results.Addresstype,
		Category:           results.Category,
		District:           results.District,
		Asn:                results.Asn,
		As:                 results.As,
	}
	return rsp, nil
}

func (g *GoIp) GetLocationInfoByIp(ip string, duration time.Duration) (*IpLocation, error) {
	cachedIpInfo, ok := gocache.Default().Get(ip)
	if ok {
		return cachedIpInfo.(*IpLocation), nil
	}
	rst, err, _ := g.singlefightGForIp.Do(ip, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		ipInfo, err := Default().GetIpLocation(ctx, ip)
		if err != nil {
			return nil, err
		}

		ipInfo.IsoCode = strings.ToLower(ipInfo.IsoCode)
		if ipInfo.IsoCode == "cn" {
			//ipInfo.Isp = ChinaISPNameToPinyin(ipInfo.Isp)
		}
		gocache.Default().SharedCache.Set(ip, ipInfo, duration)
		return ipInfo, nil
	})

	if err == nil {
		return rst.(*IpLocation), nil
	}

	return nil, err
}

// ----------------------------------------------------------------
type IpLocation struct {
	IsoCode       string `json:"iso_code"` //iso 编码 https://zh.m.wikipedia.org/zh/ISO_3166-1
	Name          string `json:"name"`
	Continent     string `json:"continent"`      //洲
	ContinentCode string `json:"continent_code"` //洲编码
	Isp           string `json:"isp"`            //isp
}

type ipLocationResp struct {
	Code int64      `json:"code"`
	Msg  string     `json:"msg"`
	Data IpLocation `json:"data"`
}

func (g *GoIp) IsLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	return ip.IsPrivate()
}

// 获取IP地址的地理位置信息 需要提供IpServiceUrl
func (g *GoIp) GetIpLocation(ctx context.Context, ip string) (*IpLocation, error) {
	netIP := net.ParseIP(ip)
	if g.IsLocalIP(netIP) {
		return &IpLocation{}, nil
	}
	request := &gohttp.Request{
		Url:     g.IpServiceUrl,
		Timeout: time.Second * 5,
	}
	request.SetQueryParams("ip", ip)
	gh := gohttp.GoHttp[ipLocationResp]{
		Request: request,
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	rst, err := gh.HttpGet(ctx)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	if rst.Code == 0 {
		return &rst.Data, nil
	}

	return nil, fmt.Errorf("request error: %w", err)
}

//---
