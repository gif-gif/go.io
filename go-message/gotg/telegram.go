package gotg

import gotgbot "github.com/gif-gif/go.io/go-tgbot"

type TelegramHook struct {
	AccessToken string
	GotgBot     *gotgbot.GoTgBot
}

// SendMessageText Function to send message
func (t *TelegramHook) SendMessageText(chatId int64, text string) error {
	_, err := t.GotgBot.SendMsgText(chatId, text)
	return err
}
