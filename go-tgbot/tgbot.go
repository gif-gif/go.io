package gotgbot

import (
	"errors"
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/util/gconv"
	"gopkg.in/telebot.v3"
	"time"
)

type GoTgBot struct {
	bot     *telebot.Bot
	reply   *telebot.ReplyMarkup
	config  *TelegramBot
	started bool
}

// 离线测试
func CreateOfflineBot(config *TelegramBot) (*GoTgBot, error) {
	bot, err := telebot.NewBot(telebot.Settings{Synchronous: true, Offline: true})
	if err != nil {
		return nil, err
	}

	return &GoTgBot{
		bot:   bot,
		reply: bot.NewMarkup(),
	}, nil
}

// 同一个产品只会存在一个
func Create(config *TelegramBot) (*GoTgBot, error) {
	if config.Product == "" {
		return nil, fmt.Errorf("product must be empty")
	}

	if config.Timeout == 0 {
		config.Timeout = 10
	}

	if config.SignTimeout == 0 {
		config.SignTimeout = 20
	}

	pref := telebot.Settings{
		Token:  config.Token,
		Poller: &telebot.LongPoller{Timeout: time.Duration(config.Timeout) * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}
	return &GoTgBot{
		bot:    bot,
		config: config,
		reply:  bot.NewMarkup(),
	}, nil
}

// /start
func (g *GoTgBot) CreateWebAppStartCommand(menuButtonText string) {
	var (
		menu = &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{
				{
					Text:   menuButtonText,
					WebApp: &telebot.WebApp{URL: g.config.WebAppUrl},
				},
			}},
		}
	)

	g.Handle("/start", func(c telebot.Context) error {
		return c.Send(g.config.StartReply, menu)
	})
}

func (g *GoTgBot) CreateMyAccountCommand(commandText string) {
	g.Handle(commandText, func(c telebot.Context) error {
		accountMd5, err := gmd5.Encrypt(gconv.String(c.Chat().ID))
		if err != nil {
			return err
		}
		_, err = g.SendMsgText(c.Chat().ID, accountMd5)
		if err != nil {
			return errors.New("failed to send message")
		}
		return nil
	})
}

func (g *GoTgBot) GetBot() *telebot.Bot {
	return g.bot
}

func (g *GoTgBot) StopBot() {
	if !g.started {
		return
	}
	g.started = false
	g.bot.Stop()
}

func (g *GoTgBot) StartBot() {
	g.started = true
	g.bot.Start()
}

func (g *GoTgBot) GetConfig() *TelegramBot {
	return g.config
}

func (g *GoTgBot) SendMsgText(chatId int64, msg string, opts ...interface{}) (*telebot.Message, error) {
	return g.Send(ChatID(chatId), msg, opts...)
}

func (g *GoTgBot) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	return g.bot.Send(to, what, opts...)
}

func (g *GoTgBot) SendAlbum(to telebot.Recipient, a telebot.Album, opts ...interface{}) ([]telebot.Message, error) {
	return g.bot.SendAlbum(to, a, opts...)
}

func (g *GoTgBot) BanSenderChat(chat *telebot.Chat, sender telebot.Recipient) error {
	return g.bot.BanSenderChat(chat, sender)
}

func (g *GoTgBot) Handle(endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) {
	g.bot.Handle(endpoint, h, m...)
}

func (g *GoTgBot) SendFromDiskPhotos(to telebot.Recipient, fromDiskFiles []string) ([]telebot.Message, error) {
	albums := telebot.Album{}
	for _, file := range fromDiskFiles {
		albums = append(albums, &telebot.Photo{File: telebot.FromDisk(file)})
	}

	return g.SendAlbum(to, albums)
}

func (g *GoTgBot) SendFromUrlPhotos(to telebot.Recipient, urls []string) ([]telebot.Message, error) {
	albums := telebot.Album{}
	for _, url := range urls {
		albums = append(albums, &telebot.Photo{File: telebot.FromURL(url)})
	}
	return g.SendAlbum(to, albums)
}

func (g *GoTgBot) SendFromDiskVideos(to telebot.Recipient, fromDiskFiles []string) ([]telebot.Message, error) {
	//v := &telebot.Video{File: telebot.FromURL("http://video.mp4")}
	albums := telebot.Album{}
	for _, file := range fromDiskFiles {
		albums = append(albums, &telebot.Video{File: telebot.FromDisk(file)})
	}
	return g.SendAlbum(to, albums)
}

func (g *GoTgBot) SendFromUrlVideos(to telebot.Recipient, urls []string) ([]telebot.Message, error) {
	//v := &telebot.Video{File: telebot.FromURL("http://video.mp4")}
	albums := telebot.Album{}
	for _, url := range urls {
		albums = append(albums, &telebot.Video{File: telebot.FromURL(url)})
	}
	return g.SendAlbum(to, albums)
}

func (g *GoTgBot) SendFromUrlAudios(to telebot.Recipient, urls []string) ([]telebot.Message, error) {
	rest := []telebot.Message{}
	for _, url := range urls {
		r, err := g.Send(to, &telebot.Audio{File: telebot.FromURL(url)})
		if err != nil {
			return nil, err
		}
		rest = append(rest, *r)
	}
	return rest, nil
}

func (g *GoTgBot) SendFromDiskAudios(to telebot.Recipient, fromDiskFiles []string) ([]telebot.Message, error) {
	rest := []telebot.Message{}
	for _, file := range fromDiskFiles {
		r, err := g.Send(to, &telebot.Audio{File: telebot.FromDisk(file)})
		if err != nil {
			return nil, err
		}
		rest = append(rest, *r)
	}
	return rest, nil
}

func (g *GoTgBot) Use(middleware ...telebot.MiddlewareFunc) {
	g.bot.Use(middleware...)
}

func (g *GoTgBot) UseMiddleware(middleware func(c telebot.Context, next telebot.HandlerFunc, parameters ...any) error, parameters ...any) {
	mm := func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			err := middleware(c, next, parameters...)
			return err
		}
	}
	g.bot.Use(mm)
}

func (g *GoTgBot) CreateMiddleware(middleware func(c telebot.Context, next telebot.HandlerFunc, parameters ...any) error, parameters ...any) telebot.MiddlewareFunc {
	mm := func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			err := middleware(c, next, parameters...)
			return err
		}
	}

	return mm
}

func (g *GoTgBot) ReplyMarkup() *telebot.ReplyMarkup {
	return g.reply
}

func (g *GoTgBot) CheckSign(ts int64, sign string) bool {
	return goutils.CheckSign(g.config.SignSecret, g.config.SignTimeout, ts, sign)
}
