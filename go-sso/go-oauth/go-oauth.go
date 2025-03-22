package gooauth

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Endpoint struct {
	TokenURL      string           `yaml:"TokenURL" json:"tokenURL,optional"`
	AuthURL       string           `yaml:"AuthURL" json:"authURL,optional"`
	DeviceAuthURL string           `yaml:"DeviceAuthURL" json:"deviceAuthURL,optional"`
	AuthStyle     oauth2.AuthStyle `yaml:"AuthStyle" json:"authStyle,optional"`
} // `yaml:"Endpoint" json:"endpoint"`

type Config struct {
	Name string `yaml:"Name" json:"name,optional"`
	// Token 信息
	AccessToken  string `yaml:"AccessToken" json:"accessToken,optional"`
	RefreshToken string `yaml:"RefreshToken" json:"refreshToken,optional"`
	ExpiresIn    int64  `yaml:"ExpiresIn" json:"expiresIn,optional"` //秒s

	// 授权参数 OAuthConfig oauth2.Config
	ClientId     string    `yaml:"ClientId" json:"clientId,optional"`
	ClientSecret string    `yaml:"ClientSecret" json:"clientSecret,optional"`
	RedirectURL  string    `yaml:"RedirectURL" json:"redirectURL,optional"`
	Endpoint     *Endpoint `yaml:"Endpoint" json:"endpoint"`
	Scopes       []string  `yaml:"Scopes" json:"scopes,optional"`
}

type GoOAuth struct {
	Config      Config
	Token       *oauth2.Token
	OAuthConfig oauth2.Config
}

func New(config Config) *GoOAuth {
	if config.Endpoint == nil { //默认google 平台
		config.Endpoint = &Endpoint{
			TokenURL:      google.Endpoint.TokenURL,
			AuthURL:       google.Endpoint.AuthURL,
			DeviceAuthURL: google.Endpoint.DeviceAuthURL,
			AuthStyle:     google.Endpoint.AuthStyle,
		}
	}

	oConfig := oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			TokenURL:      config.Endpoint.TokenURL,
			AuthURL:       config.Endpoint.AuthURL,
			DeviceAuthURL: config.Endpoint.DeviceAuthURL,
			AuthStyle:     config.Endpoint.AuthStyle,
		},
		Scopes: config.Scopes,
	}

	token := &oauth2.Token{ //兼容处理，初始化的token
		AccessToken:  config.AccessToken,
		RefreshToken: config.RefreshToken,
		ExpiresIn:    config.ExpiresIn,
	}

	return &GoOAuth{
		Token:       token,
		OAuthConfig: oConfig,
		Config:      config,
	}
}

func (c *GoOAuth) TokenSource(ctx context.Context) oauth2.TokenSource {
	return c.OAuthConfig.TokenSource(ctx, c.Token)
}

// 获取授权url
//
// 自定义参数
//
//	param1 := oauth2.SetAuthURLParam("param1", "value1")
//	param2 := oauth2.SetAuthURLParam("param2", "value2")
func (c *GoOAuth) AuthUrl(state string, opts ...oauth2.AuthCodeOption) string {
	url := c.OAuthConfig.AuthCodeURL(state, opts...)
	return url
}

// 获取token
//
// 自定义参数
//
//	param1 := oauth2.SetAuthURLParam("param1", "value1")
//	param2 := oauth2.SetAuthURLParam("param2", "value2")
func (c *GoOAuth) Exchange(ctx context.Context, authorizationCode string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	token, err := c.OAuthConfig.Exchange(ctx, authorizationCode, opts...)
	if err != nil {
		return nil, err
	}
	c.Config.AccessToken = token.AccessToken
	c.Config.RefreshToken = token.RefreshToken
	c.Token = token
	return token, nil
}

// 刷新token
func (c *GoOAuth) RefreshToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := c.OAuthConfig.TokenSource(ctx, &oauth2.Token{
		RefreshToken: c.Config.RefreshToken,
	}).Token()
	if err != nil {
		return nil, err
	}
	c.Config.AccessToken = token.AccessToken
	c.Config.RefreshToken = token.RefreshToken
	c.Token = token
	return token, nil
}

// 详细的信息，注意：只有在刷新过token的时候才会有值
func (c *GoOAuth) GetToken() *oauth2.Token {
	return c.Token
}

func (c *GoOAuth) GetAccessToken() string {
	return c.Config.AccessToken
}

func (c *GoOAuth) GetRefreshToken() string {
	return c.Config.RefreshToken
}
