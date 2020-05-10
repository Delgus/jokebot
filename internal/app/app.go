package app

import (
	"fmt"
	"io"

	"github.com/delgus/jokebot/internal/bots"
)

// Logger interface for log errors
type Logger interface {
	Error(...interface{})
}

// Options for app
type Options struct {
	EOFText               string
	InternalErrorText     string
	HelpText              string
	NotCorrectCommandText string
}

// App is my app with bot
type App struct {
	Listener
	Notifier
	bots.Bot
	Logger
	Options
}

// Run my app
func (a *App) Run(pattern string) {
	if err := a.Listen(pattern); err != nil {
		a.Error(err)
	}

	for message := range a.Message() {
		switch m := message.(type) {
		case Message:
			if m.Text == "help" {
				a.sendMessage(m.UserID, a.HelpText)
				continue
			}

			result, err := a.Command(bots.Command{
				Name: m.Text,
				Args: bots.Args{
					UserID: m.UserID,
				},
			})

			if err == io.EOF {
				a.sendMessage(m.UserID, a.EOFText)
			} else if err == bots.ErrWrongCommand {
				a.sendMessage(m.UserID, a.NotCorrectCommandText)
			} else if err != nil {
				a.Error(fmt.Errorf("bot error: %v", err))
				a.sendMessage(m.UserID, a.InternalErrorText)
			} else {
				a.sendMessage(m.UserID, result)
			}

		case ErrorMessage:
			a.Error(fmt.Errorf("listener error: %v", m.Error))
			a.sendMessage(m.UserID, a.InternalErrorText)

		default:
			a.Error(fmt.Errorf("unknown message from listener"))
		}

	}
}

func (a *App) sendMessage(userID int, result string) {
	if err := a.SendMessage(userID, result); err != nil {
		a.Error(fmt.Errorf("notifier error: %v", err))
	}
}
