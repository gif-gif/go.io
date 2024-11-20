package gometa

// API
const (
	ApiAdNetworkAnalytics = "/%s/adnetworkanalytics" // 变现广告网络分析
)

// “每分钟最多执行 250 次查询”
const limitMax = 500 //返回的行数。限制：同步请求的数量上限为 2,000。

// 常用纬度 group by 数据
// breakdowns  breakdowns=['breakdown_1', 'breakdown_2',...]
const (
	BREAKDOWN_AD_SPACE        = "ad_space"        //按广告专区细分
	BREAKDOWN_COUNTRY         = "country"         //按国家/地区细分
	BREAKDOWN_DELIVERY_METHOD = "delivery_method" //如果指标来源于通过 Audience Network 竞价投放的广告，则按 standard 或 bidding 细分。仅适用于使用变现管理工具的发行商。
	BREAKDOWN_fAIL_REASON     = "fail_reason"     //仅适用于 fb_ad_network_no_fill 和 fb_ad_network_no_bid 指标。
	BREAKDOWN_PLACEMENT       = "placement"       //按版位编号细分。不能与 placement_name 一起使用。
	BREAKDOWN_PLACEMENT_NAME  = "placement_name"  //按版位编号和名称细分。不能与 placement 一起使用。
	BREAKDOWN_PLATFORM        = "platform"        //按平台细分。可以是 ios、android、mobile_web 或 instant_games。
	BREAKDOWN_PROPERTY        = "property"        //按资产编号细分
)

// filters
const (
	FILTER_COUNTRY         = "country"         // country 以逗号分隔的双字母国家/地区缩写的清单
	FILTER_PLACEMENT       = "placement"       // placement  版位编号。限制：如果展示次数不足，值是 REDACTED。
	FILTER_DELIVERY_METHOD = "delivery_method" // delivery_method  standard 或 bidding
	FILTER_PLATFORM        = "platform"        // 可以是 ios（移动应用）、android（移动应用）、mobile_web 或 instant_games。
)

// aggregation_period=hour|day|total
const (
	AGGREGATION_PERIOD_HOUR  = "hour"
	AGGREGATION_PERIOD_DAY   = "day"
	AGGREGATION_PERIOD_TOTAL = "total"
)

// 指标
const (
	Metrics_AD_NETWORK_BIDDING_BID_RATE = "fb_ad_network_bidding_bid_rate" // 竞价响应率
	Metrics_AD_NETWORK_BIDDING_REQUEST  = "fb_ad_network_bidding_request"  // 竞价请求数量
	Metrics_AD_NETWORK_BIDDING_RESPONSE = "fb_ad_network_bidding_response" // 竞价响应数量
	Metrics_AD_NETWORK_BIDDING_WIN_RATE = "fb_ad_network_bidding_win_rate" // 竞价工具赢得竞拍的比率
	Metrics_AD_NETWORK_CLICK            = "fb_ad_network_click"            // 点击量
	Metrics_AD_NETWORK_CPM              = "fb_ad_network_cpm"              // 有效千次展示费用 (eCPM)
	Metrics_AD_NETWORK_CTR              = "fb_ad_network_ctr"              // 预估点击率
	Metrics_AD_NETWORK_FILL_RATE        = "fb_ad_network_fill_rate"        // 广告请求填充率
	Metrics_AD_NETWORK_FILLED_REQUEST   = "fb_ad_network_filled_request"   // 填充的广告请求数量
	Metrics_AD_NETWORK_IMP              = "fb_ad_network_imp"              // 展示次数
	Metrics_AD_NETWORK_NO_BID           = "fb_ad_network_no_bid"           // 无响应竞价主因数量 仅适用于用作单个指标 fail_reason 细分条件的情况
	Metrics_AD_NETWORK_NO_FILL          = "fb_ad_network_no_fill"          // 无填充主因数量仅适用于用作单个指标 fail_reason 细分条件的情况
	Metrics_AD_NETWORK_REQUEST          = "fb_ad_network_request"          // 广告请求数量
	Metrics_AD_NETWORK_REVENUE          = "fb_ad_network_revenue"          // 预估收入
	Metrics_AD_NETWORK_SHOW_RATE        = "fb_ad_network_show_rate"        // 展示数除以填充请求数
)

// ResponseData 数据
type AudienceData struct {
	QueryId string `json:"query_id"`
	Results []struct {
		Time       string `json:"time"`
		Metric     string `json:"metric"`
		Breakdowns []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"breakdowns"`
		Value string `json:"value"`
	} `json:"results"`
	OmittedResults []struct {
		Time       string `json:"time"`
		Metric     string `json:"metric"`
		Breakdowns []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"breakdowns"`
	} `json:"omitted_results"`
}

type AudienceFilter struct {
	Field    string   `json:"field"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// https://developers.facebook.com/docs/audience-network/optimization/report-api/guide-v2/
// filters=[{'field':'country', 'operator':'in', 'values':['US', 'JP']}]
type AudienceDataRequest struct {
	AggregationPeriod string `url:"aggregation_period"` //aggregation_period=hour|day|total 按 day（默认）、hour 或 total 汇总结果。限制：如要按小时汇总结果，您必须使用 since 和 until 查询至少 2 天内的结果。
	//Since 限制：
	//如要使用 Unix 时间戳，您的查询范围必须至少为 1 小时。
	//在同步请求中，您的请求范围最多为 8 天。
	//数据只会保留 540 天。如要请求的数据时间范围超过 $currentDate - 539 days，则系统不会返回更多数据。
	Since          string           `url:"since"` //since=YYYY-MM-DD 或 since=1548880485 查询的开始限制（始终包含边界值）。如果未添加此参数，默认为过去 7 天。
	Until          string           `url:"until"` //until=YYYY-MM-DD 或 until=1548880485+86400 查询的结束限制（默认不包含边界值，如果查询的汇总数据精确到小时，则包含边界值）
	Filter         []AudienceFilter `url:"filters"`
	Breakdowns     []string         `url:"breakdowns"`
	Metrics        []string         `url:"metrics"`
	Limit          int64            `url:"limit"`           //返回的行数。限制：同步请求的数量上限为 2,000。并限定每个指标获取 多少个响应
	OrderingColumn string           `url:"ordering_column"` //time|value,默认值为 time。
	OrderingType   string           `url:"ordering_type"`   //ascending|descending ,默认值为 descending。升序或降序。
	After          string           `url:"after"`           //下一页游标。
}

type AudienceDataResponse struct {
	BaseResponse[[]AudienceData]
}
