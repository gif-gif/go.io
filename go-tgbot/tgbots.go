package gotgbot

import (
	"fmt"
)

type GoTgBots struct {
	bots map[string]*GoTgBot
}

func New() *GoTgBots {
	return &GoTgBots{
		bots: make(map[string]*GoTgBot),
	}
}

// 同一个产品只会存在一个，后一个添加的覆盖前面添加
func (g *GoTgBots) CreateBot(config *TelegramBot) error {
	if config.Product == "" {
		return fmt.Errorf("product must be empty")
	}

	old := g.bots[config.Product]
	if old != nil {
		return fmt.Errorf("Already exists a bot for %s", config.Product)
	}

	bot, err := Create(config)
	if err != nil {
		return err
	}

	g.bots[config.Product] = bot
	return nil
}

func (g *GoTgBots) StopBot(product string) {
	bot := g.GetBot(product)
	bot.StopBot()
}

func (g *GoTgBots) StartBot(product string) {
	bot := g.GetBot(product)
	bot.StartBot()
}

func (g *GoTgBots) DestroyBot(product string) {
	bot := g.GetBot(product)
	bot.StopBot()
	delete(g.bots, product)
}

func (g *GoTgBots) GetBot(product string) *GoTgBot {
	return g.bots[product]
}
