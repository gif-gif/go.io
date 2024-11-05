package gometa

import "github.com/gogf/gf/util/gconv"

// 广告系列，广告组，广告 字段通用（有冗余）
const (
	accountFields = "id,account_id,name,account_status,balance,currency,business_name,business"
	allDataFields = "name,objective,status,effective_status,campaign_id,account_id,adset_id" //拉取全部数据时用的字段
	adFields      = "adset_name,account_id,campaign_id,adset_id,ad_id,impressions,cpm,cpc,spend,clicks,objective,location,cost_per_unique_click,account_name,ctr,actions"
)

type CampaignStatus string

const (
	ACTIVE   CampaignStatus = "ACTIVE"
	PAUSED   CampaignStatus = "PAUSED"
	DELETED  CampaignStatus = "DELETED"
	ARCHIVED CampaignStatus = "ARCHIVED"
)

// fb actions key
const (
	OmniActivateApp  = "omni_activate_app"
	MobileAppInstall = "mobile_app_install"
	OmniPurchase     = "omni_purchase"
	//messaging_first_reply             = "onsite_conversion.messaging_first_reply"
	//post_engagement                   = "post_engagement"
	//page_engagement                   = "page_engagement"
	//comment                           = "comment"
	//messaging_conversation_started_7d = "onsite_conversion.messaging_conversation_started_7d"
	//fb_mobile_activate_app            = "app_custom_event.fb_mobile_activate_app"
	//omni_app_install                  = "omni_app_install"
	//video_view                        = "video_view"
	//post_reaction                     = "post_reaction"
	//link_click                        = "link_click"
	//post_save                         = "onsite_conversion.post_save"
)

const (
	ApiAccount          = "/%s/client_ad_accounts"
	ApiAccountCampaigns = "/act_%s/campaigns"
	ApiAccountAdsets    = "/act_%s/adsets"
	ApiAccountAds       = "/act_%s/ads"
	ApiDataDetails      = "/%s/insights"
	ApiRefreshToken     = "/oauth/access_token"
)

type Paging struct {
	Limit  int    `json:"limit,optional"`
	Before string `json:"before,optional"`
	After  string `json:"after,optional"`
}

// 通用返回数据结构
type BaseResponse[T any] struct {
	Paging struct {
		Cursors struct {
			Before string `json:"before,optional"`
			After  string `json:"after,optional"`
		} `json:"cursors,optional"`
		Next string `json:"next,optional"`
	} `json:"paging,optional"`

	Data T `json:"data,optional"`
}

// 简要数据项（冗余：如：获取广告系列数据时 CampaignId，AdsetId 都为空 ）
type AllDataItem struct {
	Name            string `json:"name,optional"`
	Status          string `json:"status,optional"`
	EffectiveStatus string `json:"effective_status,optional"`
	CampaignId      string `json:"campaign_id,optional"`
	AccountId       string `json:"account_id,optional"`
	AdsetId         string `json:"adset_id,optional"`
	Id              string `json:"id,optional"`
}

type CampaignDetails struct {
	AccountId          string `json:"account_id,optional"`
	CampaignId         string `json:"campaign_id,optional"`
	AdsetId            string `json:"adset_id,optional"`
	AdsetName          string `json:"adset_name,optional"`
	AdId               string `json:"ad_id,optional"`
	Impressions        string `json:"impressions,optional"`
	Cpm                string `json:"cpm,optional"`
	Cpc                string `json:"cpc,optional"`
	Spend              string `json:"spend,optional"`
	Clicks             string `json:"clicks,optional"`
	Conversions        string `json:"conversions,optional"`
	Objective          string `json:"objective,optional"`
	CostPerUniqueClick string `json:"cost_per_unique_click,optional"`
	AccountName        string `json:"account_name,optional"`
	Ctr                string `json:"ctr,optional"`
	DateStart          string `json:"date_start,optional"`
	DateStop           string `json:"date_stop,optional"`
	Country            string `json:"country,optional"`

	Actions []struct {
		ActionType string `json:"action_type,optional"`
		Value      string `json:"value,optional"`
	} `json:"actions,optional"`
}

func (f *CampaignDetails) GetActionStat(actionType string) int64 {
	for _, action := range f.Actions {
		if action.ActionType == actionType {
			return gconv.Int64(action.Value)
		}
	}

	return 0
}

// facebook 投放账户信息结构
type Account struct {
	Id            string `json:"id,optional"`
	AccountId     string `json:"account_id,optional"`
	Name          string `json:"name,optional"`
	AccountStatus int    `json:"account_status,optional"`
	Balance       string `json:"balance,optional"`
	Currency      string `json:"currency,optional"`

	//BusinessName  string  `json:"business_name,optional"`
	//Business      struct {
	//	Id   string `json:"id,optional"`
	//	Name string `json:"name,optional"`
	//} `json:"business,optional"`
}

type TokenItem struct {
	AccessToken string `json:"access_token,optional"`
	TokenType   string `json:"token_type,optional"`
	ExpiresIn   int    `json:"expires_in,optional"`

	Error struct {
		Message   string `json:"message,optional"`
		Type      string `json:"type,optional"`
		Code      int    `json:"code,optional"`
		FbtraceId string `json:"fbtrace_id,optional"`
	} `json:"error,optional"`
}

type ApiRefreshTokenRequest struct {
	GrantType       string `url:"grant_type"`
	ClientId        string `url:"client_id"`
	ClientSecret    string `url:"client_secret"`
	FbExchangeToken string `url:"fb_exchange_token"`
}

type ApiAccessTokenRequest struct {
	RedirectUri  string `url:"redirect_uri"`
	ClientId     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	Code         string `url:"code,optional"`
}

type RequestData struct {
	Limit int `url:"limit"`
	//Before string `json:"before"`
	After string `url:"after"`

	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`

	DateStart string `url:"date_start"`
	DateStop  string `url:"date_stop"`
	TimeRange string `url:"time_range"` //DateStart And DateStop are here for query

	TimeIncrement int    `url:"time_increment"` // 1表示以天为单位
	Timezone      string `url:"time_zone"`      //日期时区Asia/Shanghai
	Breakdowns    string `url:"breakdowns"`     //['country'] 以国家纬度group by 查询
}

type RequestAccessKes struct {
	AccessToken string `url:"access_token"`
}

type DetailsDataRequest struct {
	RequestData
	TimeIncrement int    `url:"time_increment"`
	Timezone      string `url:"time_zone"`  //Asia/Shanghai
	Breakdowns    string `url:"breakdowns"` //['country'] 以国家纬度group by 查询
}

type AccountResponse struct {
	BaseResponse[[]Account]
}

type AllDataResponse struct {
	BaseResponse[[]AllDataItem]
}

type DataDetailResponse struct {
	BaseResponse[[]CampaignDetails]
}

type TokenResponse struct {
	TokenItem
}

// ---------meta变现端--------------------------------
type EncryptedEcpmReq struct {
	RequestId   string   `json:"request_id"`
	Ecpms       []string `json:"ecpms"`
	AccessToken string   `json:"access_token"`
	SyncApi     bool     `json:"sync_api"`
}

type EncryptedEcpmRes struct {
	RequestId string `json:"request_id"`
	Success   struct {
		Value    float64 `json:"value"`
		Accuracy string  `json:"accuracy"`
	} `json:"success"`
	Error struct {
		Reason                 string `json:"reason"`
		Description            string `json:"description"`
		NoImpressionCount      int    `json:"no_impression_count"`
		InvalidImpressionCount int    `json:"invalid_impression_count"`
	} `json:"error"`
}
