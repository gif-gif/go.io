package gotgbot

import "strconv"

type TelegramBot struct {
	Product    string `yaml:"Product" json:"product"`
	Token      string `yaml:"Token" json:"token"`
	WebAppUrl  string `yaml:"WebAppUrl,optional" json:"webAppUrl,optional"`
	StartReply string `yaml:"StartReply" json:"startReply,optional"`
	Timeout    int    `yaml:"Timeout,optional" json:"timeout,optional"` //s
	Open       bool   `yaml:"Open,optional" json:"open,optional"`       //是否开启
	Secret     string `yaml:"Secret,optional" json:"secret,optional"`   //请求密钥
}

type ChatID int64

// Recipient returns chat ID (see Recipient interface).
func (i ChatID) Recipient() string {
	return strconv.FormatInt(int64(i), 10)
}
