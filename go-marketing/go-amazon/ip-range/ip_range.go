package amazon_iprange

import (
	"context"
	gohttp "github.com/gif-gif/go.io/go-http"
)

var (
	ipRange *IpRange
)

type AmazonIp4 struct {
	IpPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type AmazonIpV6 struct {
	Ipv6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type IpRange struct {
	SyncToken    string       `json:"syncToken"`
	CreateDate   string       `json:"createDate"`
	Prefixes     []AmazonIp4  `json:"prefixes"`
	Ipv6Prefixes []AmazonIpV6 `json:"ipv6_prefixes"`
}

func LoadRangeIps() (*IpRange, error) {
	url := "https://ip-ranges.amazonaws.com/ip-ranges.json"
	req := &gohttp.Request{
		Url:     url,
		Method:  gohttp.GET,
		Headers: map[string]string{"User-Agent": "github.com/gif-gif/go.io"},
	}
	gh := gohttp.GoHttp[IpRange]{
		Request: req,
	}
	res, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func LoadRangeIpsAndCache() error {
	url := "https://ip-ranges.amazonaws.com/ip-ranges.json"
	req := &gohttp.Request{
		Url:     url,
		Method:  gohttp.GET,
		Headers: map[string]string{"User-Agent": "github.com/gif-gif/go.io"},
	}
	gh := gohttp.GoHttp[IpRange]{
		Request: req,
	}
	res, err := gh.HttpGet(context.Background())
	if err != nil {
		return err
	}
	ipRange = res
	return nil
}

func GetCacheIpRanges() *IpRange {
	return ipRange
}
