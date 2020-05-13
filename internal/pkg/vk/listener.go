package vk

import (
	"encoding/json"
	"net/http"

	"github.com/SevereCloud/vksdk/callback"
	"github.com/SevereCloud/vksdk/object"
	"github.com/delgus/jokebot/internal/pkg/app"
)

// Listener listen vk api
type Listener struct {
	cb    callback.Callback
	mChan chan interface{}
}

// NewListener return new vk listener
func NewListener(confirmToken, secretKey string) *Listener {
	cb := callback.Callback{
		ConfirmationKey: confirmToken,
		SecretKey:       secretKey,
	}
	return &Listener{
		cb:    cb,
		mChan: make(chan interface{}),
	}
}

// Message implement listener interface for app
func (l *Listener) Message() <-chan interface{} {
	return l.mChan
}

// Listen - start listening messages
func (l *Listener) Listen(pattern string) error {
	l.cb.MessageNew(func(obj object.MessageNewObject, groupID int) {
		userID := obj.Message.FromID

		command, err := getCommand(obj.Message.Payload, obj.Message.Text)
		if err != nil {
			l.mChan <- app.ErrorMessage{UserID: userID, Error: err}
			return
		}

		l.mChan <- app.Message{UserID: userID, Text: command}
	})
	http.HandleFunc(pattern, l.cb.HandleFunc)
	return nil
}

func getCommand(payload, text string) (string, error) {
	if payload == "" {
		return text, nil
	}
	payloadMap := make(map[string]string, 1)
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		return "", err
	}
	return payloadMap["command"], nil
}
