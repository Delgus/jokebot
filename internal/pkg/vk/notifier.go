package vk

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/params"
)

// Notifier - vk notifier
type Notifier struct {
	vk  *api.VK
	bsh func(*params.MessagesSendBuilder)
}

// NewNotifier return new vk notifier
func NewNotifier(accessToken string) *Notifier {
	return &Notifier{vk: api.NewVK(accessToken)}
}

// SendMessage implement app notifier
func (n *Notifier) SendMessage(userID int, text string) error {
	b := params.NewMessagesSendBuilder()
	b.PeerID(userID)
	b.RandomID(0)
	b.DontParseLinks(false)
	b.Message(text)

	if n.bsh != nil {
		n.bsh(b)
	}

	if _, err := n.vk.MessagesSend(b.Params); err != nil {
		return err
	}

	return nil
}

// SetBeforeSendHook set hook for change message before send
func (n *Notifier) SetBeforeSendHook(hook func(*params.MessagesSendBuilder)) {
	n.bsh = hook
}
