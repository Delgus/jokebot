package vkhooks

import (
	"fmt"
	"strconv"

	api "github.com/delgus/go-vk/callback-api"
	"github.com/delgus/go-vk/client"
	"github.com/delgus/jokebot/internal/app"
	"github.com/delgus/jokebot/internal/inrastructure/store/sql"
	"github.com/sirupsen/logrus"
)

const (
	jokeCommand         = "joke"
	categoryListCommand = "list"
	helpCommand         = "help"
)

var (
	jokesAreOverText         = "К сожалению шутки закончились"
	jokesCategoryAreOverText = jokesAreOverText + " Попробуйте другую категорию!"
	internalErrorText        = "К сожалению произошла ошибка. Попробуйте получить шутку позднее"
)

//OnMessageNew - hook for on message event
func OnMessageNew(jokeRepo *sql.JokeRepo, vk *client.VKClient) func(message api.VKMessageObject) {
	return func(message api.VKMessageObject) {
		switch message.Text {
		case jokeCommand:
			joke, err := jokeRepo.GetNewJoke(message.FromID)
			if err == app.ErrorJokeNotFound {
				sendMessage(vk, message.FromID, jokesAreOverText)
				return
			}
			if err != nil {
				logrus.Error(err)
				sendMessage(vk, message.FromID, internalErrorText)
				return
			}
			sendMessage(vk, message.FromID, joke.Text)

		case categoryListCommand:
			list, err := jokeRepo.GetJokeCategoryList()
			if err != nil {
				logrus.Error(err)
				sendMessage(vk, message.FromID, internalErrorText)
				return
			}
			var listMessage string
			for _, l := range list {
				listMessage = listMessage + fmt.Sprintf("%d. %s\n", l.ID, l.Name)
			}
			sendMessage(vk, message.FromID, listMessage)

		case helpCommand:
			sendMessage(vk, message.FromID, helpMessageText())

		default:
			categoryID, err := strconv.Atoi(message.Text)
			if err != nil {
				mess := "Неверная команда! \n" + helpMessageText()
				sendMessage(vk, message.FromID, mess)
				return
			}
			joke, err := jokeRepo.GetNewJokeByCategory(message.FromID, categoryID)
			if err == app.ErrorJokeNotFound {
				sendMessage(vk, message.FromID, jokesCategoryAreOverText)
				return
			}
			if err != nil {
				logrus.Error(err)
				sendMessage(vk, message.FromID, internalErrorText)
				return
			}
			sendMessage(vk, message.FromID, joke.Text)
		}
	}
}

func helpMessageText() string {
	return `
Команды для бота:

list - список категорий анекдотов

Чтобы получить анекдот из категории,отправьте номер категории

joke - возвращает анекдот из любой категории

help - помощь
`
}

func keyboard() string {
	return `
	{
  "one_time": false,
  "buttons": [
    [
      {
        "action": {
          "type": "text",
          "payload": "joke",
          "label": "Анекдот"
        },
        "color": "negative"
      },
      {
        "action": {
          "type": "text",
          "payload": "list",
          "label": "Категории анекдотов"
        },
        "color": "primary"
      },
      {
        "action": {
          "type": "text",
          "payload": "help",
          "label": "Помошь"
        },
        "color": "positive"
      }
    ]
  ]
}
	`
}

func sendMessage(vk *client.VKClient, userID int, text string) {
	if err := vk.MessagesSend(client.Message{
		UserID:   userID,
		Message:  jokesAreOverText,
		Keyboard: keyboard(),
	}); err != nil {
		logrus.Error(err)
	}
}
