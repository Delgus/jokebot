package tg

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifier tg notifier
type Notifier struct {
	client *tg.BotAPI
	bsh    func(*tg.MessageConfig)
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

	if n.bsh != nil {
		n.bsh(&msg)
	}

	if _, err := n.client.Send(msg); err != nil {
		return err
	}

	return nil
}

// SetBeforeSendHook set hook for change message before send
func (n *Notifier) SetBeforeSendHook(hook func(m *tg.MessageConfig)) {
	n.bsh = hook
}
