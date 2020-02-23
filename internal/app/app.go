package app

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
