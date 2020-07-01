package sql

import (
	"fmt"
	"io"

	"github.com/doug-martin/goqu/v8"
)

// JokeRepo realization for app.JokeRepo interface
type JokeRepo struct {
	db *goqu.Database
}

type jokeRow struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
}

type jokeCategoryRow struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// NewJokeRepo return new repository
func NewJokeRepo(db *goqu.Database) *JokeRepo {
	return &JokeRepo{db: db}
}

// GetNewJoke return new Joke
func (j *JokeRepo) GetNewJoke(userID int, client string) (string, error) {
	var row jokeRow
	found, err := j.db.From("joke").
		Select("joke.id", "text").
		LeftJoin(goqu.T("joke_user"),
			goqu.On(
				goqu.Ex{"joke_user.joke_id": goqu.I("joke.id")},
				goqu.I("user_id").Eq(userID),
				goqu.I("client").Eq(client),
			),
		).
		Where(goqu.I("user_id").IsNull()).
		ScanStruct(&row)
	if err != nil {
		return "", fmt.Errorf("infrastructure.store.sql: can't get joke: %v", err)
	}
	if !found {
		return "", io.EOF
	}
	insert := j.db.Insert("joke_user").Rows(
		goqu.Record{"user_id": userID, "client": client, "joke_id": row.ID},
	).Executor()

	if _, err := insert.Exec(); err != nil {
		return "", fmt.Errorf("infrastructure.store.sql: can't save user joke link: %v", err)
	}
	return row.Text, nil
}

// GetNewJokeByCategory return new joke by category
func (j *JokeRepo) GetNewJokeByCategory(userID, client string, categoryID int) (string, error) {
	var row jokeRow
	found, err := j.db.
		From("joke").
		Select("joke.id", "text").
		LeftJoin(goqu.T("joke_user"),
			goqu.On(
				goqu.Ex{"joke_user.joke_id": goqu.I("joke.id")},
				goqu.I("user_id").Eq(userID),
				goqu.I("client").Eq(client),
			),
		).
		Where(goqu.I("user_id").IsNull(), goqu.I("category_id").Eq(categoryID)).
		ScanStruct(&row)
	if err != nil {
		return "", fmt.Errorf("jokebot.store.sql: can't get joke: %v", err)
	}
	if !found {
		return "", io.EOF
	}
	insert := j.db.Insert("joke_user").Rows(
		goqu.Record{"user_id": userID, "client": client, "joke_id": row.ID},
	).Executor()

	if _, err := insert.Exec(); err != nil {
		return "", fmt.Errorf("jokebot.store.sql: can't save user joke link: %v", err)
	}
	return row.Text, nil
}

// GetJokeCategoryList get list of joke categories
func (j *JokeRepo) GetJokeCategoryList() (string, error) {
	var rows []jokeCategoryRow
	err := j.db.From("joke_category").
		Select(&jokeCategoryRow{}).
		ScanStructs(&rows)
	if err != nil {
		return "", fmt.Errorf("jokebot.store.sql: can't get joke category list: %v", err)
	}
	return stringList(rows), nil
}

func stringList(rows []jokeCategoryRow) string {
	var message string
	for _, r := range rows {
		message += fmt.Sprintf("%d. %s\n", r.ID, r.Name)
	}
	return message
}
