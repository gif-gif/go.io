package gomessage

import (
	"github.com/gif-gif/go.io/go-message/goding"
	"github.com/gif-gif/go.io/go-message/gofeishu"
)

var __DingDing *goding.Webhook

// 全局初始化一个DingDing 发送器对象
func InitDingDing(accessToken string, secret string) {
	__DingDing = CreateDingDing(accessToken, secret)
}

// 新创建一个DingDing 发送器对象
func CreateDingDing(accessToken string, secret string) *goding.Webhook {
	ding := &goding.Webhook{
		AccessToken: accessToken,
		Secret:      secret,
	}

	return ding
}

func DingDing(text string, at ...string) error {
	if __DingDing == nil {
		return nil
	}
	return __DingDing.SendMessageText(text, at...)
}

func FeiShu(hookUrl string, text string) error {
	return gofeishu.FeiShu(hookUrl, text)
}
