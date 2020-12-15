package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/whuwxl/twitter-media-downloader-bot/internal/pkg/service"
	"github.com/whuwxl/twitter-media-downloader-bot/internal/pkg/worker"
	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	Twitter struct {
		ConsumerKey    string `yaml:"consumerKey"`
		ConsumerSecret string `yaml:"consumerSecret"`
		AccessToken    string `yaml:"accessToken"`
		AccessSecret   string `yaml:"accessSecret"`
	} `yaml:"twitter"`
}

// NewConfig ...
func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath ...
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags ...
func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	flag.Parse()
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func main() {
	configPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	telegram, _ := service.NewTelegramService(config.Telegram.Token)
	twitter, _ := service.NewTwitterService(config.Twitter.ConsumerKey, config.Twitter.ConsumerSecret, config.Twitter.AccessToken, config.Twitter.AccessSecret)
	worker := worker.NewDownloanderWorker(telegram, twitter)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := worker.Telegram.Bot.GetUpdatesChan(u)

	for update := range updates {
		worker.FetchAndSend(update)
	}
}
