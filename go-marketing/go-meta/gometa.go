package gometa

import (
	"context"

	gohttp "github.com/gif-gif/go.io/go-http"
	"github.com/gif-gif/go.io/go-utils/gocrypto"
	"github.com/google/go-querystring/query"
)

// 1美元=1000000微分
const (
	DOLLAR_UNIT = 1000000
)

type Paging struct {
	Limit  int    `json:"limit,optional"`
	Before string `json:"before,optional"`
	After  string `json:"after,optional"`
}

// meta 通用返回数据结构
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

type Config struct {
	Name string `yaml:"Name" json:"name,optional"`
	//请求参数
	ApiVersion  string `yaml:"ApiVersion" json:"apiVersion,optional"`
	AccessToken string `yaml:"AccessToken" json:"accessToken"`

	// 授权参数
	ClientId     string `yaml:"ClientId" json:"clientId"`
	ClientSecret string `yaml:"ClientSecret" json:"clientSecret"`
	RedirectUri  string `yaml:"RedirectUri" json:"redirectUri"`

	//基础API
	baseApi               string
	currentVersionBaseApi string
}

type GoMeta struct {
	Config Config
}

func New(config Config) *GoMeta {
	mm := &GoMeta{}
	mm.UpdateConfig(config)
	return mm
}

func (m *GoMeta) UpdateConfig(config Config) {
	if config.ApiVersion == "" {
		config.ApiVersion = "v22.0"
	}
	if config.baseApi == "" {
		config.baseApi = "https://graph.facebook.com"
	}
	config.currentVersionBaseApi = config.baseApi + "/" + config.ApiVersion
	m.Config = config
}

// 刷新token接口
func (m *GoMeta) RefreshToken() (*TokenResponse, error) {
	req := &ApiRefreshTokenRequest{
		ClientId:        m.Config.ClientId,
		ClientSecret:    m.Config.ClientSecret,
		GrantType:       "fb_exchange_token",
		FbExchangeToken: m.Config.AccessToken,
	}
	api := m.Config.currentVersionBaseApi + ApiRefreshToken
	params, _ := query.Values(req)

	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[TokenResponse]{
		Request: request,
	}
	result, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取token
func (c *GoMeta) Exchange(authorizationCode string) (*TokenResponse, error) {
	req := &ApiAccessTokenRequest{
		ClientId:     c.Config.ClientId,
		ClientSecret: c.Config.ClientSecret,
		Code:         authorizationCode,
		RedirectUri:  c.Config.RedirectUri,
	}

	api := c.Config.currentVersionBaseApi + ApiRefreshToken
	params, _ := query.Values(req)

	request := &gohttp.Request{
		Url:          api,
		ParamsValues: params,
	}
	gh := gohttp.GoHttp[TokenResponse]{
		Request: request,
	}
	result, err := gh.HttpGet(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 授权URL
//
// DOC: https://developers.facebook.com/docs/marketing-api/overview/authorization
func (c *GoMeta) AuthUrl(scope string) string {
	return c.Config.currentVersionBaseApi + "/dialog/oauth?client_id=" + c.Config.ClientId + "&redirect_uri=" + gocrypto.UrlEncode(c.Config.RedirectUri) + "&scope=" + gocrypto.UrlEncode(scope)
}
