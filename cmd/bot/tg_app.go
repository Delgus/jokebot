package main

import (
	easybot "github.com/delgus/easy-bot"
	"github.com/delgus/easy-bot/clients/tg"
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func tgApp(cfg *config, bot easybot.Bot, opts easybot.Options, logger easybot.Logger) (*easybot.App, error) {
	// tg notifier
	tgNotifier, err := tg.NewNotifier(cfg.TGAccessToken)
	if err != nil {
		return nil, err
	}
	tgNotifier.SetBeforeSendHook(func(m *tba.MessageConfig) {
		m.ReplyMarkup = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton("joke"),
				tba.NewKeyboardButton("list"),
				tba.NewKeyboardButton("help"),
			),
		)
	})

	// tg listener
	tgListener, err := tg.NewListener(cfg.TGAccessToken, cfg.TGWebhook)
	if err != nil {
		return nil, err
	}

	// tg bot
	return &easybot.App{
		Notifier: tgNotifier,
		Bot:      bot,
		Listener: tgListener,
		Logger:   logger,
		Options:  opts,
	}, nil
}
