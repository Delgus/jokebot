package callback

import (
	"encoding/json"
	"net/http"

	"github.com/SevereCloud/vksdk/callback"
	"github.com/SevereCloud/vksdk/object"
)

type (
	Logger interface {
		Error(...interface{})
	}

	Service interface {
		Command(userID int, command string)
		NotifyAboutInternalError(userID int)
	}

	VKListener struct {
		service Service
		cb      callback.Callback
		logger  Logger
	}
)

func NewVKListener(confirmToken, secretKey string, service Service, logger Logger) *VKListener {
	cb := callback.Callback{
		ConfirmationKey: confirmToken,
		SecretKey:       secretKey,
	}
	return &VKListener{
		service: service,
		cb:      cb,
		logger:  logger,
	}
}

func (vkl *VKListener) Listen(pattern string) {
	vkl.cb.MessageNew(func(obj object.MessageNewObject, groupID int) {
		userID := obj.Message.FromID

		command, err := getCommand(obj.Message.Payload, obj.Message.Text)
		if err != nil {
			vkl.logger.Error(err)
			vkl.service.NotifyAboutInternalError(userID)
			return
		}

		vkl.service.Command(userID, command)
	})
	http.HandleFunc(pattern, vkl.cb.HandleFunc)
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
