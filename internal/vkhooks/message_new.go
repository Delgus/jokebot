package vkhooks

import (
	"fmt"
	"io"
	"strconv"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/params"
	"github.com/SevereCloud/vksdk/object"
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
func OnMessageNew(jokeRepo *sql.JokeRepo, vk *api.VK) func(obj object.MessageNewObject, groupID int) {
	return func(obj object.MessageNewObject, groupID int) {
		text := obj.Message.Text
		userID := obj.Message.FromID

		switch text {
		case jokeCommand:
			joke, err := jokeRepo.GetNewJoke(obj.Message.FromID)
			if err == io.EOF {
				if _, err := vk.MessagesSend(message(userID, jokesAreOverText)); err != nil {
					logrus.Error(err)
				}
				return
			}
			if err != nil {
				logrus.Error(err)
				if _, err := vk.MessagesSend(message(userID, internalErrorText)); err != nil {
					logrus.Error(err)
				}
				return
			}

			if _, err := vk.MessagesSend(message(userID, joke.Text)); err != nil {
				logrus.Error(err)
			}

		case categoryListCommand:
			list, err := jokeRepo.GetJokeCategoryList()
			if err != nil {
				logrus.Error(err)
				if _, err := vk.MessagesSend(message(userID, internalErrorText)); err != nil {
					logrus.Error(err)
				}
				return
			}

			var listMessage string
			for _, l := range list {
				listMessage = listMessage + fmt.Sprintf("%d. %s\n", l.ID, l.Name)
			}

			if _, err := vk.MessagesSend(message(userID, listMessage)); err != nil {
				logrus.Error(err)
			}

		case helpCommand:
			if _, err := vk.MessagesSend(message(userID, helpMessageText())); err != nil {
				logrus.Error(err)
			}

		default:
			categoryID, err := strconv.Atoi(text)
			if err != nil {
				m := "Неверная команда! \n" + helpMessageText()
				if _, err := vk.MessagesSend(message(userID, m)); err != nil {
					logrus.Error(err)
				}
				return
			}

			joke, err := jokeRepo.GetNewJokeByCategory(userID, categoryID)
			if err == io.EOF {
				if _, err := vk.MessagesSend(message(userID, jokesCategoryAreOverText)); err != nil {
					logrus.Error(err)
				}
				return
			}
			if err != nil {
				logrus.Error(err)
				if _, err := vk.MessagesSend(message(userID, internalErrorText)); err != nil {
					logrus.Error(err)
				}
				return
			}

			if _, err := vk.MessagesSend(message(userID, joke.Text)); err != nil {
				logrus.Error(err)
			}
		}
	}
}

func message(userID int, text string) api.Params {
	b := params.NewMessagesSendBuilder()
	b.PeerID(userID)
	b.DontParseLinks(false)
	b.Message(text)
	return b.Params
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
