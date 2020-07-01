package jokebot

import (
	"strconv"

	easybot "github.com/delgus/easy-bot"
)

type (
	// JokeBot is my bot
	JokeBot struct {
		r JokeRepo
	}

	// JokeRepo is repository for storage of jokes
	JokeRepo interface {
		GetNewJoke(userID int, client string) (string, error)
		GetNewJokeByCategory(userID int, client string, categoryID int) (string, error)
		GetJokeCategoryList() (string, error)
	}
)

const (
	// JokeCommand - joke command
	JokeCommand = "joke"
	// CategoryListCommand - category list command
	CategoryListCommand = "list"
)

// NewBot return new bot
func NewBot(r JokeRepo) *JokeBot {
	return &JokeBot{r: r}
}

// Command implement interface Bot
func (j *JokeBot) Command(command easybot.Command) (string, error) {
	switch command.Name {
	case JokeCommand:
		return j.r.GetNewJoke(command.Args.UserID, string(command.Args.Client))

	case CategoryListCommand:
		return j.r.GetJokeCategoryList()

	default:
		cat, err := strconv.Atoi(command.Name)
		if err != nil {
			return "", easybot.ErrWrongCommand
		}
		return j.r.GetNewJokeByCategory(command.Args.UserID, string(command.Args.Client), cat)
	}
}
