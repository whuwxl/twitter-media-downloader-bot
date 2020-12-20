package worker

import (
	"net/url"
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
	if update.ChannelPost == nil && update.Message == nil {
		return nil
	}

	var text string
	var chatID int64
	var messageID int

	if update.ChannelPost != nil {
		text = update.ChannelPost.Text
		chatID = update.ChannelPost.Chat.ID
		messageID = update.ChannelPost.MessageID
	}
	if update.Message != nil {
		text = update.Message.Text
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
	}

	url, err := url.Parse(text)
	if err != nil {
		return err
	}
	s := strings.Split(url.Path, "/")
	id, err := strconv.ParseInt(s[len(s)-1], 10, 64)
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
			} else if m.Type == "video" {
				maxBitrate := 0
				url := ""
				for _, v := range m.VideoInfo.Variants {
					if v.Bitrate > maxBitrate {
						maxBitrate = v.Bitrate
						url = v.URL
					}
				}
				urls = append(urls, url)
			}
		}
	}

	if len(urls) == 1 {
		if err := worker.Telegram.SendDocument(chatID, urls[0]); err != nil {
			worker.Telegram.SendMessage(chatID, err.Error()+"\n"+strings.Join(urls, "\n"), messageID)
			return err
		}
	} else if len(urls) > 1 {
		if err := worker.Telegram.SendDocuments(chatID, urls); err != nil {
			worker.Telegram.SendMessage(chatID, err.Error()+"\n"+strings.Join(urls, "\n"), messageID)
			return err
		}
	} else {
		return nil
	}

	worker.Telegram.DeleteMessage(chatID, messageID)

	return nil
}
