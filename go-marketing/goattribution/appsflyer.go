package goattribution

import (
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

/*
{
    "af_tranid": "XW2mV9SveRKJ0p-mNip_hA",
    "af_c_id": "ss_ads4eachs_android_HappyFruit2048_br_1208",
    "af_adset_id": "1806913189",
    "pid": "mintegral_int",
    "af_prt": "ads4eachs",
    "af_adset": "icon_512x512",
    "af_ad": "icon_512x512",
    "af_siteid": "mtg1145003732",
    "af_ad_id": "1806913189",
    "c": "ads4eachs_android_HappyFruit2048_br_1208"
}
*/

// type AppsFlyer struct {
// 	AfTranID  string `json:"af_tranid"`
// 	AfCID     string `json:"af_c_id"`
// 	AfAdsetID string `json:"af_adset_id"`
// 	Pid       string `json:"pid"`
// 	AfPRT     string `json:"af_prt"`
// 	AfAdset   string `json:"af_adset"`
// 	AfAd      string `json:"af_ad"`
// 	AfSiteID  string `json:"af_siteid"`
// 	AfAdID    string `json:"af_ad_id"`
// 	C         string `json:"c"`
// }

type AppsFlyerAttributeHandler struct {
}

func (h *AppsFlyerAttributeHandler) Channel() string {
	return "appsflyer"
}

func (h *AppsFlyerAttributeHandler) Match(queryParams url.Values) bool {
	return len(queryParams.Get("af_tranid")) > 0
}

func (h *AppsFlyerAttributeHandler) Handle(queryParams url.Values) (*AttributeInfo, error) {
	info, err := CreateAttributeInfo(queryParams.Get("af_c_id"), queryParams.Get("c"))
	if err != nil {
		logx.Errorf("AppsFlyerAttributeHandler handle %v queryParams:%+v", err, queryParams)
	}
	// info.Channel = userdef.CHANNEL_APPSFLYER
	info.Channel = info.CampaignChannel
	info.CampaignName += "_" + queryParams.Get("af_adset")
	return info, nil
}
