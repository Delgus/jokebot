package vk

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/params"
)

// Notifier - vk notifier
type Notifier struct {
	vk *api.VK
}

// NewNotifier return new vk notifier
func NewNotifier(accessToken string) *Notifier {
	return &Notifier{
		vk: api.Init(accessToken),
	}
}

// SendMessage implement app notifier
func (n *Notifier) SendMessage(userID int, text string) error {
	b := params.NewMessagesSendBuilder()
	b.PeerID(userID)
	b.RandomID(0)
	b.DontParseLinks(false)
	b.Message(text)
	b.Keyboard(`{
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

	if _, err := n.vk.MessagesSend(b.Params); err != nil {
		return err
	}

	return nil
}
