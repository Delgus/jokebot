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

	VKCallback struct {
		service Service
		cb      callback.Callback
		logger  Logger
	}
)

func NewVKCallback(confirmToken, secretKey string, service Service, logger Logger) *VKCallback {
	cb := callback.Callback{
		ConfirmationKey: confirmToken,
		SecretKey:       secretKey,
	}
	cb.MessageNew(func(obj object.MessageNewObject, groupID int) {
		userID := obj.Message.FromID

		command, err := getCommand(obj.Message.Payload, obj.Message.Text)
		if err != nil {
			logger.Error(err)
			service.NotifyAboutInternalError(userID)
			return
		}

		service.Command(userID, command)
	})
	return &VKCallback{
		service: service,
		cb:      cb,
		logger:  logger,
	}
}

func (c *VKCallback) HandleFunc(w http.ResponseWriter, r *http.Request) {
	c.cb.HandleFunc(w, r)
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
