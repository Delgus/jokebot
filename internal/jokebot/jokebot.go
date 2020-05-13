package jokebot

import (
	"strconv"

	"github.com/delgus/jokebot/internal/pkg/app"
)

type (
	// JokeBot is my bot
	JokeBot struct {
		r JokeRepo
	}

	// JokeRepo is repository for storage of jokes
	JokeRepo interface {
		GetNewJoke(userID int) (string, error)
		GetNewJokeByCategory(userID int, categoryID int) (string, error)
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
func (j *JokeBot) Command(command app.Command) (string, error) {
	switch command.Name {
	case JokeCommand:
		return j.r.GetNewJoke(command.Args.UserID)

	case CategoryListCommand:
		return j.r.GetJokeCategoryList()

	default:
		cat, err := strconv.Atoi(command.Name)
		if err != nil {
			return "", app.ErrWrongCommand
		}
		return j.r.GetNewJokeByCategory(command.Args.UserID, cat)
	}

}
