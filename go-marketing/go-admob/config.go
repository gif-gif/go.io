package goadmob

type Config struct {
	Name         string `yaml:"Name" json:"name,optional"`
	AccountId    string `yaml:"AccountId" json:"accountId,optional"`
	ClientId     string `yaml:"ClientId" json:"clientId"`
	ClientSecret string `yaml:"ClientSecret" json:"clientSecret"`
	RedirectUrl  string `yaml:"RedirectUrl" json:"redirectUrl"`
	State        string `yaml:"State" json:"state"`

	AccessToken  string `yaml:"AccessToken" json:"accessToken"`
	RefreshToken string `yaml:"RefreshToken" json:"refreshToken"`
}
