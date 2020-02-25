package app

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sirupsen/logrus"
)

type JokeCategory struct {
	ID   int
	Name string
}

type Joke struct {
	ID   int
	Text string
}

type JokeRepo interface {
	GetNewJoke(userID int) (*Joke, error)
	GetNewJokeByCategory(userID int, categoryID int) (*Joke, error)
	GetJokeCategoryList() ([]JokeCategory, error)
}

type Notifier interface {
	SendMessage(int, string)
}

type JokeService struct {
	JokeCommand string
	ListCommand string
	HelpCommand string

	JokesAreOverText       string
	TryAnotherCategoryText string
	InternalErrorText      string
	HelpMessageText        string
	NotCorrectCommandText  string

	Notifier Notifier
	Repo     JokeRepo
}

func (j *JokeService) Command(userID int, command string) {
	switch command {
	case j.JokeCommand:
		joke, err := j.Repo.GetNewJoke(userID)
		if err == io.EOF {
			j.Notifier.SendMessage(userID, j.JokesAreOverText)
			return
		}
		if err != nil {
			logrus.Error(err)
			j.Notifier.SendMessage(userID, j.InternalErrorText)
			return
		}

		j.Notifier.SendMessage(userID, joke.Text)

	case j.ListCommand:
		list, err := j.Repo.GetJokeCategoryList()
		if err != nil {
			logrus.Error(err)
			j.Notifier.SendMessage(userID, j.InternalErrorText)
			return
		}

		j.Notifier.SendMessage(userID, listMessage(list))

	case j.HelpCommand:
		j.Notifier.SendMessage(userID, j.HelpMessageText)

	default:
		categoryID, err := strconv.Atoi(command)
		if err != nil {
			j.Notifier.SendMessage(userID, j.NotCorrectCommandText+j.HelpMessageText)
			return
		}

		joke, err := j.Repo.GetNewJokeByCategory(userID, categoryID)
		if err == io.EOF {
			j.Notifier.SendMessage(userID, j.JokesAreOverText+" "+j.TryAnotherCategoryText)
			return
		}
		if err != nil {
			logrus.Error(err)
			j.Notifier.SendMessage(userID, j.InternalErrorText)
			return
		}

		j.Notifier.SendMessage(userID, joke.Text)
	}
}

func listMessage(list []JokeCategory) string {
	var message string
	for _, l := range list {
		message += fmt.Sprintf("%d. %s\n", l.ID, l.Name)
	}
	return message
}
