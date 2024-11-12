package gooauth

import (
	"context"
	"golang.org/x/oauth2"
)

// 授权接口
type IGoOAuth interface {
	AuthUrl() string
	Exchange(ctx context.Context, authorizationCode string) (*oauth2.Token, error)
	RefreshToken(ctx context.Context) (*oauth2.Token, error)
	GetToken() *oauth2.Token
}
