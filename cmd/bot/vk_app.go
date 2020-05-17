package main

import (
	"github.com/SevereCloud/vksdk/api/params"
	easybot "github.com/delgus/easy-bot"
	"github.com/delgus/easy-bot/clients/vk"
)

func vkApp(cfg *config, bot easybot.Bot, opts easybot.Options, logger easybot.Logger) *easybot.App {
	vkNotifier := vk.NewNotifier(cfg.VKAccessToken)
	vkNotifier.SetBeforeSendHook(func(m *params.MessagesSendBuilder) {
		m.Keyboard(`{
			"buttons": [
			  [
				{
				  "action": {
					"type": "text",
					"label": "Анекдот",
					"payload": "{\"command\":\"joke\"}"
				  },
				  "color": "positive"
				}
			  ],
			  [
				{
				  "action": {
					"type": "text",
					"label": "Категории анекдотов",
					"payload": "{\"command\":\"list\"}"
				  },
				  "color": "negative"
				}
			  ],
			  [
				{
				  "action": {
					"type": "text",
					"label": "Помощь",
					"payload": "{\"command\":\"help\"}"
				  },
				  "color": "primary"
				}
			  ]
			]
		  }`)
	})
	return &easybot.App{
		Notifier: vkNotifier,
		Bot:      bot,
		Listener: vk.NewListener(cfg.VKConfirmToken, cfg.VKSecretKey),
		Logger:   logger,
		Options:  opts,
	}
}
