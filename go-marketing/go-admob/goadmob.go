package goadmob

import (
	"context"
	"errors"
	gohttp "github.com/gif-gif/go.io/go-http"
	golog "github.com/gif-gif/go.io/go-log"
	gooauth "github.com/gif-gif/go.io/go-sso/go-oauth"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/oauth2"
	"google.golang.org/api/admob/v1"
	"google.golang.org/api/option"
	"math"
	"time"
)

const (
	DefaultCurrencyCode = "USD"
	DefaultLanguageCode = "en-US"
)

// 常用纬度
const (
	FILTER_APP              = "APP"
	FILTER_COUNTRY          = "COUNTRY"
	FILTER_DATE             = "DATE"
	FILTER_PLATFORM         = "PLATFORM"
	FILTER_AD_UNIT          = "AD_UNIT"
	FILTER_FORMAT           = "FORMAT"
	FILTER_APP_VERSION_NAME = "APP_VERSION_NAME"
)

const (
	AD_REQUESTS        = "AD_REQUESTS"
	CLICKS             = "CLICKS"
	ESTIMATED_EARNINGS = "ESTIMATED_EARNINGS"
	IMPRESSIONS        = "IMPRESSIONS"
	IMPRESSION_CTR     = "IMPRESSION_CTR"
	IMPRESSION_RPM     = "IMPRESSION_RPM"
	MATCHED_REQUESTS   = "MATCHED_REQUESTS"
	MATCH_RATE         = "MATCH_RATE"
	SHOW_RATE          = "SHOW_RATE"
)

var DefaultMetrics = []string{AD_REQUESTS, CLICKS, ESTIMATED_EARNINGS, IMPRESSIONS, IMPRESSION_CTR, IMPRESSION_RPM, MATCHED_REQUESTS, MATCH_RATE, SHOW_RATE}

type ReportReq struct {
	Dimensions      []string //查询维度列表
	AdUnits         []string //广告位
	Countries       []string //国家
	AdFormats       []string //原生、横幅、插屏、开屏、激励视频
	Platforms       []string //应用的移动操作系统平台（例如“Android”或“iOS”）。
	AppVersionNames []string //对于 Android，应用版本名称可以在 PackageInfo 中的 versionName 中找到。对于 iOS，可以在 CFBundleShortVersionString 中找到应用版本名称。警告：该维度与 ESTIMATED_EARNINGS 和 OBSERVED_ECPM 指标不兼容。
	MaxReportRows   int64    //最大返回数量
	Metrics         []string //查询字段、注意：这里的查询字段在不同的纬度下有互斥情况

	StartDate admob.Date
	EndDate   admob.Date

	CurrencyCode string //default currency USD
	LanguageCode string //default language en-US

	RetryWaitTime time.Duration
	RetryCount    int

	//Date Month Week
}

type ResponseItem struct {
	Date            string
	AdUnit          string
	Country         string
	AdRequest       int64
	Clicks          int64
	Earnings        int64 //美分
	Impressions     int64
	ImpressionCtr   float64
	ImpressionRpm   int64 //美分
	MatchedRequests int64
	MatchRate       float64
	ShowRate        float64
}

// accessToken 会在60分钟后过期
type GoAdmob struct {
	ctx            context.Context
	Config         Config
	GoAuth         *gooauth.GoOAuth
	AdmobService   *admob.Service
	RequestTimeout int64 //default request timeout 30s
}

// 每次调用时都需要调用这个方法
func New(ctx context.Context, config Config) (*GoAdmob, error) {
	if len(config.AuthConfig.Scopes) <= 0 {
		config.AuthConfig.Scopes = []string{"https://www.googleapis.com/auth/admob.readonly", "https://www.googleapis.com/auth/admob.report"}
	}
	err := gooauth.Init(config.AuthConfig)
	if err != nil {
		return nil, err
	}
	o := &GoAdmob{
		ctx:    ctx,
		Config: config,
		GoAuth: gooauth.GetClient(config.AuthConfig.Name),
	}

	return o, nil
}

// RefreshToken and AdmobService admon token 有效期为60分钟，所有每次请求数据刷新下token
func (o *GoAdmob) Refresh() error {
	err := o.RefreshToken()
	if err != nil {
		return err
	}

	admobService, err := admob.NewService(o.ctx, option.WithTokenSource(o.GoAuth.TokenSource(o.ctx)))
	if err != nil {
		return err
	}

	o.AdmobService = admobService
	return nil
}

// 授权admobURL OAuth2.0
//
// 1、https://developers.google.com/admob/api/v1/getting-started?hl=zh-cn 创建一个新的授权客户端 获取client_id 和 client_secret
//
// 2、浏览器中执行 AuthUrl 方法获取授权code
//
// 3、获取code后执行 Exchange 方法获取token
func (c *GoAdmob) AuthUrl() string {
	return c.GoAuth.AuthUrl()
}

// 获取token
func (c *GoAdmob) Exchange(ctx context.Context, authorizationCode string) (*oauth2.Token, error) {
	token, err := c.GoAuth.Exchange(ctx, authorizationCode)
	if err != nil {
		golog.WithTag("goadmob").Error("token:" + err.Error())
		return nil, err
	}

	return token, nil
}

// 刷新token
func (c *GoAdmob) RefreshToken() error {
	_, err := c.GoAuth.RefreshToken(c.ctx)
	if err != nil {
		golog.WithTag("goadmob").Error("token:" + err.Error())
		return err
	}
	return nil
}

// 获取广告报表信息 这个接口不支持查询维度，返回接口有bug
//func (c *GoAdmob) GetReport() (*admob.GenerateNetworkReportResponse, error) {
//	res, err := c.AdmobService.Accounts.NetworkReport.Generate("accounts/"+c.Config.AccountId, &admob.GenerateNetworkReportRequest{
//		ReportSpec: &admob.NetworkReportSpec{
//			DateRange: &admob.DateRange{
//				EndDate: &admob.Date{
//					Day:   21,
//					Month: 8,
//					Year:  2024,
//				},
//				StartDate: &admob.Date{
//					Day:   20,
//					Month: 8,
//					Year:  2024,
//				},
//			},
//			Dimensions: []string{"DATE", "APP", "COUNTRY"},
//			//DimensionFilters: []*admob.NetworkReportSpecDimensionFilter{
//			//	{
//			//		Dimension: "COUNTRY",
//			//		MatchesAny: &admob.StringList{
//			//			Values: []string{"US"},
//			//		},
//			//	},
//			//},
//			MaxReportRows: 10,
//			Metrics:       []string{"CLICKS", "ESTIMATED_EARNINGS"},
//		},
//	}).Do()
//	return res, err
//}

// dimensions 指定查询维度，如：[]string{"DATE", "APP", "COUNTRY"}
//
// SELECT DATE, APP, COUNTRY, CLICKS, ESTIMATED_EARNINGS
// FROM NETWORK_REPORT
// WHERE DATE >= '2021-09-01' AND DATE <= '2021-09-30'
//
//	AND COUNTRY IN ('US', 'CN')
//
// GROUP BY DATE, APP, COUNTRY
// ORDER BY APP ASC, CLICKS DESC;
func (c *GoAdmob) GetReport(req *ReportReq) ([]*ResponseItem, error) {
	if req.MaxReportRows == 0 {
		return nil, errors.New("MaxReportRows is empty")
	}

	url := "/v1/accounts/" + c.Config.AccountId + "/networkReport:generate"
	dimensionFilters := []*admob.NetworkReportSpecDimensionFilter{}
	if len(req.AdUnits) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: FILTER_AD_UNIT,
			MatchesAny: &admob.StringList{
				Values: req.AdUnits,
			},
		})
	}

	if len(req.AdFormats) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: FILTER_FORMAT,
			MatchesAny: &admob.StringList{
				Values: req.AdFormats,
			},
		})
	}

	if len(req.Platforms) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: FILTER_PLATFORM,
			MatchesAny: &admob.StringList{
				Values: req.Platforms,
			},
		})
	}

	if len(req.AppVersionNames) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: FILTER_APP_VERSION_NAME,
			MatchesAny: &admob.StringList{
				Values: req.AppVersionNames,
			},
		})
	}
	if len(req.Countries) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: FILTER_COUNTRY,
			MatchesAny: &admob.StringList{
				Values: req.Countries,
			},
		})
	}

	if len(req.Metrics) == 0 {
		req.Metrics = DefaultMetrics
	}

	if req.CurrencyCode == "" {
		req.CurrencyCode = DefaultCurrencyCode
	}

	if req.LanguageCode == "" {
		req.LanguageCode = DefaultLanguageCode
	}

	params := &admob.GenerateNetworkReportRequest{
		ReportSpec: &admob.NetworkReportSpec{
			DateRange: &admob.DateRange{
				EndDate:   &req.EndDate,
				StartDate: &req.StartDate,
			},
			Dimensions:       req.Dimensions, //[]string{"DATE", "APP", "COUNTRY"}, //group by
			DimensionFilters: dimensionFilters,
			MaxReportRows:    req.MaxReportRows,
			Metrics:          req.Metrics,
			LocalizationSettings: &admob.LocalizationSettings{
				CurrencyCode: req.CurrencyCode,
				LanguageCode: req.LanguageCode,
			},
		},
	}

	dataReq := &gohttp.Request{
		Url:           url,
		RetryCount:    req.RetryCount,
		RetryWaitTime: req.RetryWaitTime,
		Timeout:       time.Second * 30,
		Body:          params,
	}

	gh := &gohttp.GoHttp[[]*admob.GenerateNetworkReportResponse]{
		Request: dataReq,
		BaseUrl: "https://admob.googleapis.com",
		Headers: map[string]string{
			"X-Google-AuthUser": "0",
			"Authorization":     "Bearer " + c.GoAuth.GetAccessToken(),
		},
	}

	res, err := gh.HttpPostJson(c.ctx)
	if err != nil {
		return nil, err
	}

	list := []*ResponseItem{}
	for _, response := range *res {
		item := ResponseItem{}
		row := response.Row
		if row != nil {
			for fieldName, value := range row.DimensionValues {
				switch fieldName {
				case FILTER_DATE:
					item.Date = value.Value
					break
				case FILTER_APP:
					item.AdUnit = value.Value
					break
				case FILTER_COUNTRY:
					item.Country = value.Value
					break
				}
			}

			for fieldName, value := range row.MetricValues {
				switch fieldName {
				case AD_REQUESTS:
					item.AdRequest = value.IntegerValue
					break
				case CLICKS:
					item.Clicks = value.IntegerValue
					break
				case ESTIMATED_EARNINGS:
					item.Earnings = value.MicrosValue / 10000
					break
				case IMPRESSIONS:
					item.Impressions = value.IntegerValue
					break
				case IMPRESSION_CTR:
					item.ImpressionCtr = value.DoubleValue
					break
				case IMPRESSION_RPM:
					item.ImpressionRpm = gconv.Int64(math.Floor(value.DoubleValue * 100))
					break
				case MATCHED_REQUESTS:
					item.MatchedRequests = value.IntegerValue
					break
				case MATCH_RATE:
					item.MatchRate = value.DoubleValue
					break
				case SHOW_RATE:
					item.ShowRate = value.DoubleValue
					break
				}
			}
			list = append(list, &item)
		}
	}

	return list, nil
}

// 获取账号下所有APP信息
func (c *GoAdmob) GetApps() (*admob.ListAppsResponse, error) {
	if c.Config.AccountId == "" {
		return nil, errors.New("accountId is empty")
	}
	res, err := c.AdmobService.Accounts.Apps.List("accounts/" + c.Config.AccountId).Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 获取当前appId下所有的广告信息
func (c *GoAdmob) GetAdUnits() (*admob.ListAdUnitsResponse, error) {
	if c.Config.AccountId == "" {
		return nil, errors.New("accountId is empty")
	}
	res, err := c.AdmobService.Accounts.AdUnits.List("accounts/" + c.Config.AccountId).Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}
