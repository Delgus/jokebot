package notify

import tg "github.com/go-telegram-bot-api/telegram-bot-api"

type (
	TGNotifier struct {
		client *tg.BotAPI
		logger Logger
	}
)

func NewTGNotifier(accessToken string, logger Logger) (*TGNotifier, error) {
	client, err := tg.NewBotAPI(accessToken)
	if err != nil {
		return nil, err
	}
	return &TGNotifier{
		client: client,
		logger: logger,
	}, nil
}

func (tgn *TGNotifier) SendMessage(userID int, text string) {
	msg := tg.NewMessage(int64(userID), text)
	_, err := tgn.client.Send(msg)
	if err != nil {
		tgn.logger.Error(err)
	}
}
