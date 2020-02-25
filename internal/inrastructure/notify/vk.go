package notify

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/params"
	"github.com/sirupsen/logrus"
)

type VKNotifier struct {
	vk *api.VK
}

func NewVKNotifier(vk *api.VK) *VKNotifier {
	return &VKNotifier{vk: vk}
}

func (n *VKNotifier) SendMessage(userID int, text string) {
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
		logrus.Error(err)
	}
}
