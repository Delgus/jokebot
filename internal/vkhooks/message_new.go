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
				if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: jokesAreOverText}); err != nil {
					logrus.Error(err)
				}
				return
			}
			if err != nil {
				logrus.Error(err)
				if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: internalErrorText}); err != nil {
					logrus.Error(err)
				}
				return
			}
			if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: joke.Text}); err != nil {
				logrus.Error(err)
			}
		case categoryListCommand:
			list, err := jokeRepo.GetJokeCategoryList()
			if err != nil {
				logrus.Error(err)
				if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: internalErrorText}); err != nil {
					logrus.Error(err)
				}
				return
			}
			var listMessage string
			for _, l := range list {
				listMessage = listMessage + fmt.Sprintf("%d. %s\n", l.ID, l.Name)
			}
			if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: listMessage}); err != nil {
				logrus.Error(err)
			}
		case helpCommand:
			if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: helpMessageText()}); err != nil {
				logrus.Error(err)
			}
		default:
			categoryID, err := strconv.Atoi(message.Text)
			if err != nil {
				mess := "Неверная команда! \n" + helpMessageText()
				if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: mess}); err != nil {
					logrus.Error(err)
				}
				return
			}
			joke, err := jokeRepo.GetNewJokeByCategory(message.FromID, categoryID)
			if err == app.ErrorJokeNotFound {
				if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: jokesCategoryAreOverText}); err != nil {
					logrus.Error(err)
				}
				return
			}
			if err != nil {
				logrus.Error(err)
				if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: internalErrorText}); err != nil {
					logrus.Error(err)
				}
				return
			}
			if err := vk.MessagesSend(client.Message{UserID: message.FromID, Message: joke.Text}); err != nil {
				logrus.Error(err)
			}
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
