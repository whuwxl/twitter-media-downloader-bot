package service

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Twitter ...
type Twitter struct {
	Client *twitter.Client
}

// NewTwitterService ...
func NewTwitterService(consumerKey string, consumerSecret string, accessToken string, accessSecret string) (*Twitter, error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return &Twitter{Client: client}, nil
}

// ShowTweet ...
func (t *Twitter) ShowTweet(id int64) (*twitter.Tweet, error) {
	params := &twitter.StatusShowParams{TrimUser: twitter.Bool(true)}
	tweet, _, err := t.Client.Statuses.Show(id, params)
	return tweet, err
}
