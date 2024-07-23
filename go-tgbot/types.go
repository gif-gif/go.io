package gotgbot

import "strconv"

type TelegramBot struct {
	Product    string `yaml:"Product" json:"product"`
	Token      string `yaml:"Token" json:"token"`
	WebAppUrl  string `yaml:"WebAppUrl,optional" json:"webAppUrl,optional"`
	StartReply string `yaml:"StartReply" json:"startReply,optional"`
	Timeout    int64  `yaml:"Timeout,optional" json:"timeout,optional"` //s

	// Product configuration (可选)
	Open         bool   `yaml:"Open,optional" json:"open,optional"`                 //是否开启
	SignSecret   string `yaml:"SignSecret,optional" json:"signSecret,optional"`     //请求签名密钥
	SignTimeout  int64  `yaml:"SignTimeout,optional" json:"signTimeout,optional"`   //请求签名过期时间-秒
	AccessSecret string `yaml:"AccessSecret,optional" json:"accessSecret,optional"` //jwt signature secret
	AccessExpire int64  `yaml:"AccessExpire,optional" json:"accessExpire,optional"` // jwt signature accessExpire
}

type ChatID int64

// Recipient returns chat ID (see Recipient interface).
func (i ChatID) Recipient() string {
	return strconv.FormatInt(int64(i), 10)
}
