package goattribution

import (
	"errors"
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"net/url"
	"strings"
)

const (
	CHANNEL_ORGANIC   = "organic"
	CHANNEL_META      = "meta"
	CHANNEL_GOOGLE    = "google"
	CHANNEL_APPSFLYER = "appsflyer"
	CHANNEL_MINTEGRAL = "mintegral"
	CHANNEL_BIGO      = "bigo"
)

type AttributeHandler interface {
	Match(queryParams url.Values) bool
	Handle(queryParams url.Values) (*AttributeInfo, error)
	Channel() string //开发者定义的渠道标识
}

type AttributeInfo struct {
	Channel         string
	CampaignId      string
	CampaignName    string
	AdCostMode      string
	CampaignChannel string
	CampaignPartner string
	UtmCampaignId   string
	UtmCampaignName string
	UtmSource       string
	UtmMedium       string
	UtmContent      string
}

func CreateAttributeInfo(queryParams url.Values, campaignId, campaignName string) (*AttributeInfo, error) {
	result := &AttributeInfo{
		UtmSource:    queryParams.Get("utm_source"),
		UtmMedium:    queryParams.Get("utm_medium"),
		UtmContent:   queryParams.Get("utm_content"),
		CampaignId:   campaignId,
		CampaignName: campaignName,
	}

	items := strings.Split(campaignName, "_")
	if len(items) > 0 {
		result.CampaignPartner = items[0]
	}
	if len(items) > 2 {
		result.AdCostMode = items[2]
	}
	if len(items) > 3 {
		result.CampaignChannel = items[3]
	}

	var err error
	if len(items) < 4 {
		err = fmt.Errorf("campaignName error:%s", campaignName)
	}
	return result, err
}

type Config struct {
	Name        string
	DecryptKeys map[string]string // 每个平台 解密key -> [平台标识]=[解密key]
}

var __clients = map[string]*AttributeManager{}

// 可以一次初始化多个Redis实例或者 多次调用初始化多个实例
func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("goattribution client [" + name + "] already exists")
		}

		__clients[name] = New(conf)
	}

	return
}

func GetClient(names ...string) *AttributeManager {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := __clients[name]; ok {
		return cli
	}
	return nil
}

func DelClient(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			delete(__clients, name)
		}
	}
}

func Default() *AttributeManager {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("goattribution").Error("no default goattribution client")

	return nil
}
