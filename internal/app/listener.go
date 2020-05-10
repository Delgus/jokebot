package app

// Message - message from client
type Message struct {
	UserID int
	Text   string
}

// ErrorMessage - message about internal error
type ErrorMessage struct {
	UserID int
	Error  error
}

// HelpMessage - message about help
type HelpMessage struct {
	UserID int
	Help   bool
}

// Listener interface for listen client
type Listener interface {
	Message() <-chan interface{}
	Listen(string) error
}
