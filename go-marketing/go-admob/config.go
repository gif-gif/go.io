package goadmob

import gooauth "github.com/gif-gif/go.io/go-sso/go-oauth"

type Config struct {
	Name       string `yaml:"Name" json:"name,optional"`
	AccountId  string `yaml:"AccountId" json:"accountId,optional"`
	AuthConfig gooauth.Config
	//AccessToken  string `yaml:"AccessToken" json:"accessToken"`
	//RefreshToken string `yaml:"RefreshToken" json:"refreshToken"`
	//ClientId     string `yaml:"ClientId" json:"clientId"`
	//ClientSecret string `yaml:"ClientSecret" json:"clientSecret"`
	//RedirectUrl  string `yaml:"RedirectUrl" json:"redirectUrl"`
	//State        string `yaml:"State" json:"state"`
}
