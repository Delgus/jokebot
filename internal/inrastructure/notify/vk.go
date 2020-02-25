package notify

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/params"
)

type (
	Logger interface {
		Error(...interface{})
	}

	VKNotifier struct {
		vk       *api.VK
		logger   Logger
		keyboard string
	}
)

func NewVKNotifier(accessToken string, logger Logger) *VKNotifier {
	return &VKNotifier{
		vk:       api.Init(accessToken),
		logger:   logger,
		keyboard: "{}",
	}
}

func (n *VKNotifier) Keyboard(keyboard string) {
	n.keyboard = keyboard
}

func (n *VKNotifier) SendMessage(userID int, text string) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(userID)
	b.RandomID(0)
	b.DontParseLinks(false)
	b.Message(text)
	b.Keyboard(n.keyboard)
	if _, err := n.vk.MessagesSend(b.Params); err != nil {
		n.logger.Error(err)
	}
}
