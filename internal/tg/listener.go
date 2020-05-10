package tg

import (
	"github.com/delgus/jokebot/internal/app"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Listener listener for telegram
type Listener struct {
	host   string
	client *tg.BotAPI
	mChan  chan interface{}
}

// Message implement listener interface for app
func (l *Listener) Message() <-chan interface{} {
	return l.mChan
}

// NewListener return new tg listener
func NewListener(accessToken, host string) (*Listener, error) {
	client, err := tg.NewBotAPI(accessToken)
	if err != nil {
		return nil, err
	}

	return &Listener{
		host:   host,
		client: client,
		mChan:  make(chan interface{}),
	}, nil
}

// Listen tg webhook
func (l *Listener) Listen(pattern string) error {
	// tg need in set webhook
	_, err := l.client.SetWebhook(tg.NewWebhook(l.host + pattern))
	if err != nil {
		return err
	}

	updates := l.client.ListenForWebhook(pattern)

	go func() {
		for update := range updates {
			if update.Message == nil || update.Message.From == nil {
				continue
			}
			userID := update.Message.Chat.ID
			command := update.Message.Text

			if command == "help" {
				l.mChan <- app.HelpMessage{UserID: int(userID), Help: true}
			}
			if err != nil {
				l.mChan <- app.ErrorMessage{UserID: int(userID), Error: err}
				continue
			}

			l.mChan <- app.Message{UserID: int(userID), Text: command}
		}
	}()
	return nil
}
