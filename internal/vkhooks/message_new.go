package vkhooks

import (
	"encoding/json"

	"github.com/SevereCloud/vksdk/object"
	"github.com/delgus/jokebot/internal/app"
	"github.com/sirupsen/logrus"
)

//OnMessageNew - hook for on message event
func OnMessageNew(service *app.JokeService) func(obj object.MessageNewObject, groupID int) {
	return func(obj object.MessageNewObject, groupID int) {
		userID := obj.Message.FromID

		command, err := getCommand(obj.Message.Payload, obj.Message.Text)
		if err != nil {
			logrus.Error(err)
			service.Notifier.SendMessage(userID, service.InternalErrorText)
			return
		}

		service.Command(userID, command)
	}
}

func getCommand(payload string, text string) (string, error) {
	if payload == "" {
		return text, nil
	}
	payloadMap := make(map[string]string, 1)
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		return "", err
	}
	return payloadMap["command"], nil
}
