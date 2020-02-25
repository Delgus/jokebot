package callback

import tg "github.com/go-telegram-bot-api/telegram-bot-api"

type (
	TGListener struct {
		host    string
		service Service
		client  *tg.BotAPI
		logger  Logger
	}
)

func NewTGListener(accessToken, host string, service Service, logger Logger) (*TGListener, error) {
	client, err := tg.NewBotAPI(accessToken)
	if err != nil {
		return nil, err
	}

	return &TGListener{
		host:    host,
		service: service,
		client:  client,
		logger:  logger,
	}, nil
}

func (tgl *TGListener) Listen(pattern string) error {
	// tg need in set webhook
	_, err := tgl.client.SetWebhook(tg.NewWebhook(tgl.host + pattern))
	if err != nil {
		return err
	}

	updates := tgl.client.ListenForWebhook(pattern)

	go func() {
		for update := range updates {
			if update.Message == nil || update.Message.From == nil {
				continue
			}
			userID := update.Message.Chat.ID

			command := update.Message.Text
			if err != nil {
				tgl.logger.Error(err)
				tgl.service.NotifyAboutInternalError(int(userID))

				continue
			}

			tgl.service.Command(int(userID), command)
		}
	}()
	return nil
}
