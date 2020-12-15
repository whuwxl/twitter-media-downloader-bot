package worker

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/whuwxl/twitter-media-downloader-bot/internal/pkg/service"
)

// DownloanderWorker ...
type DownloanderWorker struct {
	Telegram *service.Telegram
	Twitter  *service.Twitter
}

// NewDownloanderWorker ...
func NewDownloanderWorker(telegram *service.Telegram, twitter *service.Twitter) *DownloanderWorker {
	return &DownloanderWorker{Telegram: telegram, Twitter: twitter}
}

// FetchAndSend ...
func (worker *DownloanderWorker) FetchAndSend(update tgbotapi.Update) error {
	if update.ChannelPost == nil {
		return nil
	}

	s := update.ChannelPost.Text
	ss := strings.Split(s, "/")
	id, err := strconv.ParseInt(ss[len(ss)-1], 10, 64)
	if err != nil {
		return err
	}

	tweet, err := worker.Twitter.ShowTweet(id)
	if err != nil {
		return err
	}

	urls := make([]string, 0)

	if tweet.ExtendedEntities != nil && tweet.ExtendedEntities.Media != nil {
		for _, m := range tweet.ExtendedEntities.Media {
			if m.Type == "photo" {
				urls = append(urls, m.MediaURLHttps+":orig")
			}
		}
	}

	if len(urls) == 1 {
		worker.Telegram.SendDocument(update.ChannelPost.Chat.ID, urls[0])
	} else if len(urls) > 1 {
		worker.Telegram.SendDocuments(update.ChannelPost.Chat.ID, urls)
	}

	worker.Telegram.DeleteMessage(update.ChannelPost.Chat.ID, update.ChannelPost.MessageID)

	return nil
}
