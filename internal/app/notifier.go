package app

// Notifier interface notifier for clients
type Notifier interface {
	SendMessage(int, string) error
}
