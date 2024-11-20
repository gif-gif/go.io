package gometa

import (
	"context"
	gohttp "github.com/gif-gif/go.io/go-http"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/google/go-querystring/query"
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
	mm := &GoMeta{
		Config: config,
	}
	if mm.Config.ApiVersion == "" {
		mm.Config.ApiVersion = "v21.0"
	}
	if mm.Config.baseApi == "" {
		mm.Config.baseApi = "https://graph.facebook.com"
	}

	mm.Config.currentVersionBaseApi = mm.Config.baseApi + "/" + mm.Config.ApiVersion
	return mm
}

// 刷新token接口
func (m *GoMeta) RefreshToken(clientId string, clientSecret string) (*TokenResponse, error) {
	req := &ApiRefreshTokenRequest{
		ClientId:        clientId,
		ClientSecret:    clientSecret,
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
	return c.Config.currentVersionBaseApi + "/dialog/oauth?client_id=" + c.Config.ClientId + "&redirect_uri=" + goutils.UrlEncode(c.Config.RedirectUri) + "&scope=" + goutils.UrlEncode(scope)
}
