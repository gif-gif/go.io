package main

import (
	golog "github.com/gif-gif/go.io/go-log"
	gotgbot "github.com/gif-gif/go.io/go-tgbot"
	"github.com/gogf/gf/util/gconv"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"time"
)

type recipientString string

func (r recipientString) Recipient() string {
	return "5562314141"
}

type ChatID int64

// Recipient returns chat ID (see Recipient interface).
func (i ChatID) Recipient() string {
	return strconv.FormatInt(int64(i), 10)
}

func main() {
	pref := &gotgbot.TelegramBot{
		Product:    "test",
		Token:      "7107568224:AAFgdiEsDqtFvBBScIfWku9IB8jr9Dpl-dw",
		WebAppUrl:  "https://www.google.com",
		StartReply: "hello",
	}

	gobot, err := gotgbot.CreateOfflineBot(pref)
	if err != nil {
		golog.Error(err.Error())
		return
	}
	var (
		// Universal markup builders.
		menu = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{{
				{
					Text:   "This is a web app",
					WebApp: &tele.WebApp{URL: "https://google.com"},
				},
			}},
		}
	)

	gobot.Handle("/start", func(c tele.Context) error {
		return c.Send(pref.StartReply, menu)
	})

	var trace []string

	handler := func(name string) tele.HandlerFunc {
		return func(c tele.Context) error {
			trace = append(trace, name)
			return nil
		}
	}

	middleware := func(c tele.Context, next tele.HandlerFunc, params ...any) error {
		trace = append(trace, gconv.String(params[0])+":in")
		err := next(c)
		trace = append(trace, gconv.String(params[0])+":out")
		return err
	}

	createMiddleware1 := gobot.CreateMiddleware(middleware, "handler1a")
	createMiddleware2 := gobot.CreateMiddleware(middleware, "handler2a")
	createGroup1Middleware1 := gobot.CreateMiddleware(middleware, "group1")
	createGroup2Middleware2 := gobot.CreateMiddleware(middleware, "group2")
	createGroup2Middleware3 := gobot.CreateMiddleware(middleware, "handler1b")

	b := gobot.GetBot()
	gobot.UseMiddleware(middleware, "test")
	gobot.UseMiddleware(middleware, "global1")
	gobot.UseMiddleware(middleware, "global2")

	b.Handle("/a", handler("/a"), createMiddleware1, createMiddleware2)

	group := b.Group()
	group.Use(createGroup1Middleware1, createGroup2Middleware2)
	group.Handle("/b", handler("/b"), createGroup2Middleware3)

	b.ProcessUpdate(tele.Update{
		Message: &tele.Message{Text: "/a"},
	})

	golog.WithTag("gotgbot").Info(trace)
	trace = trace[:0]
	b.ProcessUpdate(tele.Update{
		Message: &tele.Message{Text: "/b"},
	})
	golog.WithTag("gotgbot").Info(trace)

	// ReplyMarkup is a part of SendOptions,
	// but often it's the only option you need
	gobot.SendMsgText(5562314141, "text")
	gobot.SendFromUrlPhotos(ChatID(5562314141), []string{"https://developer.android.google.cn/static/studio/images/run/adb_wifi-quick_settings.png?hl=zh-cn"})

	// flags: no notification && no web link preview
	go gobot.StartBot()

	time.Sleep(2 * time.Second)

	time.Sleep(1000 * time.Second)
}
