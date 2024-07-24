package main

import (
	golog "github.com/gif-gif/go.io/go-log"
	gotgbot "github.com/gif-gif/go.io/go-tgbot"
	"time"
)

// Recipient returns chat ID (see Recipient interface).

func main() {
	var chatId int64 = 5562314141
	pref := &gotgbot.TelegramBot{
		Product:    "test",
		Token:      "7107568224:AAFgdiEsDqtFvBBScIfWku9IB8jr9Dpl-dw",
		WebAppUrl:  "https://www.google.com",
		StartReply: "hello",
		ApiUrl:     "https://tgapi.goio.dev",
	}

	gobot, err := gotgbot.Create(pref) // TODO离线模式测试
	if err != nil {
		golog.Error(err.Error())
		return
	}
	gobot.CreateMyAccountCommand("/account")

	var trace []string

	//handler := func(name string) tele.HandlerFunc {
	//	return func(c tele.Context) error {
	//		trace = append(trace, name)
	//		return nil
	//	}
	//}
	//
	//middleware := func(c tele.Context, next tele.HandlerFunc, params ...any) error {
	//	trace = append(trace, gconv.String(params[0])+":in")
	//	err := next(c)
	//	trace = append(trace, gconv.String(params[0])+":out")
	//	return err
	//}
	//
	//createMiddleware1 := gobot.CreateMiddleware(middleware, "handler1a")
	//createMiddleware2 := gobot.CreateMiddleware(middleware, "handler2a")
	//createGroup1Middleware1 := gobot.CreateMiddleware(middleware, "group1")
	//createGroup2Middleware2 := gobot.CreateMiddleware(middleware, "group2")
	//createGroup2Middleware3 := gobot.CreateMiddleware(middleware, "handler1b")

	//b := gobot.GetBot()
	//gobot.UseMiddleware(middleware, "test")
	//gobot.UseMiddleware(middleware, "global1")
	//gobot.UseMiddleware(middleware, "global2")

	//b.Handle("/a", handler("/a"), createMiddleware1, createMiddleware2)
	//
	//group := b.Group()
	//group.Use(createGroup1Middleware1, createGroup2Middleware2)
	//group.Handle("/b", handler("/b"), createGroup2Middleware3)

	//b.ProcessUpdate(tele.Update{
	//	Message: &tele.Message{Text: "/a"},
	//})

	golog.WithTag("gotgbot").Info(trace)
	trace = trace[:0]
	//b.ProcessUpdate(tele.Update{
	//	Message: &tele.Message{Text: "/b"},
	//})
	golog.WithTag("gotgbot").Info(trace)

	gobot.SendMsgText(chatId, "text")
	//gobot.SendFromUrlPhotos(gotgbot.ChatID(chatId), []string{
	//	"https://developer.android.google.cn/static/studio/images/run/adb_wifi-quick_settings.png?hl=zh-cn",
	//	"https://biuvip.com/bot/how_agents_buy_tokens.png",
	//})

	go gobot.StartBot()
	//gobot.StopBot()
	time.Sleep(1000 * time.Second)
}
