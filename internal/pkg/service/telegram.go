package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Telegram ...
type Telegram struct {
	Bot *tgbotapi.BotAPI
}

// NewTelegramService ...
func NewTelegramService(token string) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &Telegram{Bot: bot}, nil
}

// DeleteMessage ...
func (t *Telegram) DeleteMessage(chatID int64, messageID int) error {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// SendDocument ...
func (t *Telegram) SendDocument(chatID int64, url string) error {
	msg := tgbotapi.NewDocumentShare(chatID, url)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// SendDocuments ...
func (t *Telegram) SendDocuments(chatID int64, urls []string) error {
	files := make([]interface{}, len(urls))
	for i, v := range urls {
		files[i] = tgbotapi.InputMediaPhoto{
			Type:  "document",
			Media: v,
		}
	}

	msg := tgbotapi.NewMediaGroup(chatID, files)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
