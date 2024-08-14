package gomessage

import (
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/go-message/goding"
	"github.com/gif-gif/go.io/go-message/gofeishu"
	"github.com/gif-gif/go.io/go-message/gotg"
	gotgbot "github.com/gif-gif/go.io/go-tgbot"
)

var __DingDing *goding.Webhook
var __Telegram *gotg.TelegramHook

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

// telegram

// Telegram 发送器对象, inChina true时走代理IP
func InitTelegram(accessToken string, inChina bool) {
	api := "https://api.telegram.com"
	if inChina {
		api = "https://tgapi.goio.dev"
	}
	pref := &gotgbot.TelegramBot{
		Product:    "goio",
		Token:      accessToken,
		StartReply: "",
		ApiUrl:     api,
	}

	gobot, err := gotgbot.Create(pref)
	if err != nil {
		golog.WithTag("gomessage").Error(err.Error())
		return
	}
	go gobot.StartBot()

	__Telegram = &gotg.TelegramHook{
		AccessToken: accessToken,
		GotgBot:     gobot,
	}
}

// chatId 个人ID或群组ID text 消息内容
func Telegram(chatId int64, text string) error {
	return __Telegram.SendMessageText(chatId, text)
}
