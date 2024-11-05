package gooauth

import (
	"context"
	"golang.org/x/oauth2"
)

type Config struct {
	Name         string `yaml:"Name" json:"name,optional"`
	AccessToken  string `yaml:"AccessToken" json:"accessToken"`
	RefreshToken string `yaml:"RefreshToken" json:"refreshToken"`
	State        string `yaml:"State" json:"state"`
	AuthConfig   oauth2.Config
}

type GoOAuth struct {
	Config Config
	Token  *oauth2.Token
}

func New(config Config) *GoOAuth {
	return &GoOAuth{
		Config: config,
	}
}

func (c *GoOAuth) TokenSource(ctx context.Context) oauth2.TokenSource {
	return c.Config.AuthConfig.TokenSource(ctx, c.Token)
}

// 获取授权url
func (c *GoOAuth) AuthUrl() string {
	url := c.Config.AuthConfig.AuthCodeURL(c.Config.State)
	return url
}

// 获取token
func (c *GoOAuth) Exchange(ctx context.Context, authorizationCode string) (*oauth2.Token, error) {
	token, err := c.Config.AuthConfig.Exchange(ctx, authorizationCode)
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
	token, err := c.Config.AuthConfig.TokenSource(ctx, &oauth2.Token{
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
