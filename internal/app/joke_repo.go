package app

import "errors"

var (
	ErrorJokeNotFound = errors.New("joke not found")
)

type JokeRepo interface {
	GetNewJoke(userID int) (*Joke, error)
	GetNewJokeByCategory(userID int, categoryID int) (*Joke, error)
	GetJokeCategoryList() ([]JokeCategory, error)
}
