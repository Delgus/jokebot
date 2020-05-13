package app

import "fmt"

// ErrWrongCommand - bot return this error when command is not correct
var ErrWrongCommand = fmt.Errorf(`incorrect command for bot`)

// Command for bot
type Command struct {
	Name string
	Args Args
}

// Args arguments for command
type Args struct {
	UserID int
}

// Bot get command and return result
type Bot interface {
	Command(Command) (string, error)
}
