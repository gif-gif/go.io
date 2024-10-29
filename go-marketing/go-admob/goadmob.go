package goadmob

import (
	"context"
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/admob/v1"
	"google.golang.org/api/option"
)

// accessToken 会在60分钟后过期
var client_id = "273488495628-h81gn5a7q5a6632j4nbapfpu6cq3454l.apps.googleusercontent.com"
var client_secret = "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"

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
func (c *GoAdmob) GetReport() (*admob.GenerateNetworkReportResponse, error) {
	res, err := c.AdmobService.Accounts.NetworkReport.Generate("accounts/pub-2200607303212623", &admob.GenerateNetworkReportRequest{
		ReportSpec: &admob.NetworkReportSpec{
			DateRange: &admob.DateRange{
				EndDate: &admob.Date{
					Day:   21,
					Month: 8,
					Year:  2021,
				},
				StartDate: &admob.Date{
					Day:   20,
					Month: 8,
					Year:  2021,
				},
			},
			Dimensions:    []string{"DATE", "APP", "COUNTRY"},
			MaxReportRows: 100000,
			Metrics:       []string{"CLICKS", "ESTIMATED_EARNINGS"},
		},
	}).Do()
	return res, err
}

// 获取账号下所有APP信息
func (c *GoAdmob) GetApps(accountId string) (*admob.ListAppsResponse, error) {
	if accountId == "" {
		return nil, errors.New("accountId is empty")
	}
	res, err := c.AdmobService.Accounts.Apps.List("accounts/" + accountId).Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 获取当前appId下所有的广告信息
func (c *GoAdmob) GetAdUnits(accountId string) (*admob.ListAdUnitsResponse, error) {
	if accountId == "" {
		return nil, errors.New("accountId is empty")
	}
	res, err := c.AdmobService.Accounts.AdUnits.List("accounts/" + accountId).Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}
