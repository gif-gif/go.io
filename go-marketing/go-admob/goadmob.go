package goadmob

import (
	"context"
	"errors"
	gohttp "github.com/gif-gif/go.io/go-http"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/oauth2"
	"google.golang.org/api/admob/v1"
	"google.golang.org/api/option"
	"math"
	"time"
)

type ReportReq struct {
	Dimensions      []string //查询维度列表
	AdUnits         []string //广告位
	AdFormats       []string //原生、横幅、插屏、开屏、激励视频
	Platforms       []string //应用的移动操作系统平台（例如“Android”或“iOS”）。
	AppVersionNames []string //对于 Android，应用版本名称可以在 PackageInfo 中的 versionName 中找到。对于 iOS，可以在 CFBundleShortVersionString 中找到应用版本名称。警告：该维度与 ESTIMATED_EARNINGS 和 OBSERVED_ECPM 指标不兼容。
	MaxReportRows   int64    //最大返回数量
	Metrics         []string //查询字段、注意：这里的查询字段在不同的纬度下有互斥情况

	StartDate admob.Date
	EndDate   admob.Date

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

var metrics = []string{"AD_REQUESTS", "CLICKS", "ESTIMATED_EARNINGS", "IMPRESSIONS", "IMPRESSION_CTR", "IMPRESSION_RPM", "MATCHED_REQUESTS", "MATCH_RATE", "SHOW_RATE"}

// accessToken 会在60分钟后过期
type GoAdmob struct {
	ctx          context.Context
	Config       Config
	AuthConfig   oauth2.Config
	AdmobService *admob.Service
	Token        *oauth2.Token
}

// 每次调用时都需要调用这个方法
func New(ctx context.Context, config Config) (*GoAdmob, error) {
	authConfig := oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	o := &GoAdmob{
		ctx:        ctx,
		AuthConfig: authConfig,
		Config:     config,
	}

	return o, nil
}

// RefreshToken and AdmobService admon token 有效期为60分钟，所有每次请求数据刷新下token
func (o *GoAdmob) Refresh() error {
	err := o.RefreshToken()
	if err != nil {
		return err
	}

	admobService, err := admob.NewService(o.ctx, option.WithTokenSource(o.AuthConfig.TokenSource(o.ctx, o.Token)))
	if err != nil {
		return err
	}

	o.AdmobService = admobService
	return nil
}

// 授权admobURL
func (c *GoAdmob) AuthUrl() string {
	url := `https://accounts.google.com/o/oauth2/v2/auth?client_id=` + c.Config.ClientId + `&redirect_uri=` + goutils.UrlEncode(c.Config.RedirectUrl) + `&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fadmob.report&prompt=consent&state=` + c.Config.State + `&response_type=code&access_type=offline`
	return url
}

// 获取token
func (c *GoAdmob) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.AuthConfig.Exchange(ctx, code)
	if err != nil {
		golog.WithTag("goadmob").Error("token:" + err.Error())
		return nil, err
	}

	return token, nil
}

// 刷新token
func (c *GoAdmob) RefreshToken() error {
	token, err := c.AuthConfig.TokenSource(c.ctx, &oauth2.Token{
		RefreshToken: c.Config.RefreshToken,
	}).Token()
	if err != nil {
		golog.WithTag("goadmob").Error("token:" + err.Error())
		return err
	}
	c.Token = token
	return nil
}

// 获取广告报表信息
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
			Dimension: "AD_UNIT",
			MatchesAny: &admob.StringList{
				Values: req.AdUnits,
			},
		})
	}

	if len(req.AdFormats) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: "FORMAT",
			MatchesAny: &admob.StringList{
				Values: req.AdFormats,
			},
		})
	}

	if len(req.Platforms) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: "PLATFORM",
			MatchesAny: &admob.StringList{
				Values: req.Platforms,
			},
		})
	}

	if len(req.AppVersionNames) > 0 {
		dimensionFilters = append(dimensionFilters, &admob.NetworkReportSpecDimensionFilter{
			Dimension: "APP_VERSION_NAME",
			MatchesAny: &admob.StringList{
				Values: req.AppVersionNames,
			},
		})
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
				CurrencyCode: "USD",
				LanguageCode: "en-US",
			},
		},
	}

	dataReq := &gohttp.Request{
		Url:     url,
		Timeout: time.Second * 20,
		Body:    params,
	}

	gh := &gohttp.GoHttp[[]*admob.GenerateNetworkReportResponse]{
		Request: dataReq,
		BaseUrl: "https://admob.googleapis.com",
		Headers: map[string]string{
			"X-Google-AuthUser": "0",
			"Authorization":     "Bearer " + c.Token.AccessToken,
		},
	}

	res, err := gh.HttpPostJson(c.ctx)
	if err != nil {
		return nil, err
	}

	//type ResponseItem struct {
	//	Date            string
	//	AdUnit          string
	//	Country         string
	//	AdRequest       int64
	//	Clicks          int64
	//	Earnings        float64
	//	Impressions     int64
	//	ImpressionCtr   float64
	//	ImpressionRpm   float64
	//	MatchedRequests int64
	//	MatchRate       float64
	//	ShowRate        float64
	//}
	//
	//var metrics = []string{"AD_REQUESTS", "CLICKS", "ESTIMATED_EARNINGS", "IMPRESSIONS", "IMPRESSION_CTR", "IMPRESSION_RPM", "MATCHED_REQUESTS", "MATCH_RATE", "SHOW_RATE"}
	list := []*ResponseItem{}
	for _, response := range *res {
		item := ResponseItem{}
		row := response.Row
		if row != nil {
			for fieldName, value := range row.DimensionValues {
				switch fieldName {
				case "DATE":
					item.Date = value.Value
					break
				case "APP":
					item.AdUnit = value.Value
					break
				case "COUNTRY":
					item.Country = value.Value
					break
				}
			}

			for fieldName, value := range row.MetricValues {
				switch fieldName {
				case "AD_REQUESTS":
					item.AdRequest = value.IntegerValue
					break
				case "CLICKS":
					item.Clicks = value.IntegerValue
					break
				case "ESTIMATED_EARNINGS":
					item.Earnings = value.MicrosValue / 10000
					break
				case "IMPRESSIONS":
					item.Impressions = value.IntegerValue
					break
				case "IMPRESSION_CTR":
					item.ImpressionCtr = value.DoubleValue
					break
				case "IMPRESSION_RPM":
					item.ImpressionRpm = gconv.Int64(math.Floor(value.DoubleValue * 100))
					break
				case "MATCHED_REQUESTS":
					item.MatchedRequests = value.IntegerValue
					break
				case "MATCH_RATE":
					item.MatchRate = value.DoubleValue
					break
				case "SHOW_RATE":
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
