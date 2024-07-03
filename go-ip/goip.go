package goip

import (
	"github.com/ip2location/ip2location-go/v9"
	"github.com/oschwald/geoip2-golang"
	"net"
)

type GoIp struct {
	dbReader   *geoip2.Reader
	dbLocation *ip2location.DB
}

func (g *GoIp) Init(config Config) error {
	db, err := geoip2.Open(config.Mmdb)
	if err != nil {
		return err
	}
	g.dbReader = db

	dbLocation, err := ip2location.OpenDB(config.Ip2locationDB)
	if err != nil {
		return err
	}
	g.dbLocation = dbLocation
	return nil
}

// 查询 IP相关信息,返回IP所属国家中文名称
func (g *GoIp) QueryDbReaderCountryForZhName(ipStr string) (*IpCountry, error) {
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
