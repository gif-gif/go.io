package gotgbot

import (
	"gopkg.in/telebot.v3"
	"strconv"
)

type TelegramBot struct {
	Product     string            `yaml:"Product" json:"product"`
	ApiUrl      string            `yaml:"ApiUrl" json:"apiUrl,optional"`
	ParseMode   telebot.ParseMode `yaml:"ParseMode" json:"parseMode,optional"`
	Token       string            `yaml:"Token" json:"token"`
	Timeout     int64             `yaml:"Timeout" json:"timeout,optional"` //s tg bot api timeout
	Offline     bool              `yaml:"Offline" json:"offline,optional"`
	Synchronous bool              `yaml:"Synchronous" json:"synchronous,optional"`

	// Product configuration (可选)
	Open         bool   `yaml:"Open" json:"open,optional"` //是否开启
	WebAppUrl    string `yaml:"WebAppUrl" json:"webAppUrl,optional"`
	StartReply   string `yaml:"StartReply" json:"startReply,optional"`
	SignSecret   string `yaml:"SignSecret" json:"signSecret,optional"`     //请求签名密钥
	SignTimeout  int64  `yaml:"SignTimeout" json:"signTimeout,optional"`   //请求签名过期时间-秒
	AccessSecret string `yaml:"AccessSecret" json:"accessSecret,optional"` //jwt signature secret
	AccessExpire int64  `yaml:"AccessExpire" json:"accessExpire,optional"` // jwt signature accessExpire
}

type ChatID int64

// Recipient returns chat ID (see Recipient interface).
func (i ChatID) Recipient() string {
	return strconv.FormatInt(int64(i), 10)
}
