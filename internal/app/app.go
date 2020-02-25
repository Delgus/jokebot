package app

import (
	"fmt"
	"io"
	"strconv"
)

type (
	JokeCategory struct {
		ID   int
		Name string
	}

	Joke struct {
		ID   int
		Text string
	}

	JokeRepo interface {
		GetNewJoke(userID int) (*Joke, error)
		GetNewJokeByCategory(userID int, categoryID int) (*Joke, error)
		GetJokeCategoryList() ([]JokeCategory, error)
	}

	Notifier interface {
		SendMessage(int, string)
	}

	Logger interface {
		Error(...interface{})
	}

	Options struct {
		JokeCommand string
		ListCommand string
		HelpCommand string

		JokesAreOverText       string
		TryAnotherCategoryText string
		InternalErrorText      string
		HelpMessageText        string
		NotCorrectCommandText  string
	}

	JokeService struct {
		*Options
		n      Notifier
		r      JokeRepo
		logger Logger
	}
)

func NewJokeService(n Notifier, r JokeRepo, l Logger, o *Options) *JokeService {
	return &JokeService{
		Options: o,
		n:       n,
		r:       r,
		logger:  l,
	}
}

func (j *JokeService) Command(userID int, command string) {
	switch command {
	case j.JokeCommand:
		joke, err := j.r.GetNewJoke(userID)
		if err == io.EOF {
			j.n.SendMessage(userID, j.JokesAreOverText)
			return
		}
		if err != nil {
			j.logger.Error(err)
			j.n.SendMessage(userID, j.InternalErrorText)
			return
		}

		j.n.SendMessage(userID, joke.Text)

	case j.ListCommand:
		list, err := j.r.GetJokeCategoryList()
		if err != nil {
			j.logger.Error(err)
			j.n.SendMessage(userID, j.InternalErrorText)
			return
		}

		j.n.SendMessage(userID, listMessage(list))

	case j.HelpCommand:
		j.n.SendMessage(userID, j.HelpMessageText)

	default:
		categoryID, err := strconv.Atoi(command)
		if err != nil {
			j.n.SendMessage(userID, j.NotCorrectCommandText+j.HelpMessageText)
			return
		}

		joke, err := j.r.GetNewJokeByCategory(userID, categoryID)
		if err == io.EOF {
			j.n.SendMessage(userID, j.JokesAreOverText+" "+j.TryAnotherCategoryText)
			return
		}
		if err != nil {
			j.logger.Error(err)
			j.n.SendMessage(userID, j.InternalErrorText)
			return
		}

		j.n.SendMessage(userID, joke.Text)
	}
}

func (j *JokeService) NotifyAboutInternalError(userID int) {
	j.n.SendMessage(userID, j.InternalErrorText)
}

func listMessage(list []JokeCategory) string {
	var message string
	for _, l := range list {
		message += fmt.Sprintf("%d. %s\n", l.ID, l.Name)
	}
	return message
}
