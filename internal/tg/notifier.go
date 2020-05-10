package tg

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifier tg notifier
type Notifier struct {
	client *tg.BotAPI
}

// NewNotifier return new tg notifier
func NewNotifier(accessToken string) (*Notifier, error) {
	client, err := tg.NewBotAPI(accessToken)
	if err != nil {
		return nil, err
	}

	return &Notifier{client: client}, nil
}

// SendMessage implement notifier interface
func (n *Notifier) SendMessage(userID int, text string) error {

	msg := tg.NewMessage(int64(userID), text)
	msg.ReplyMarkup = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("joke"),
			tg.NewKeyboardButton("list"),
			tg.NewKeyboardButton("help"),
		),
	)

	if _, err := n.client.Send(msg); err != nil {
		return err
	}

	return nil
}
