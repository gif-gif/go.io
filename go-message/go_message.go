package gomessage

import (
	"github.com/gif-gif/go.io/go-message/ding"
	"github.com/gif-gif/go.io/go-message/feishu"
)

var __DingDing *ding.Webhook

func InitDingDing(accessToken string, secret string) {
	__DingDing = &ding.Webhook{
		AccessToken: accessToken,
		Secret:      secret,
	}
}

func FeiShu(hookUrl string, text string) error {
	return feishu.FeiShu(hookUrl, text)
}

func DingDing(text string) error {
	if __DingDing == nil {
		return nil
	}
	return __DingDing.SendMessageText(text)
}

func GetDingDing(hookUrl string, text string) *ding.Webhook {
	return __DingDing
}
