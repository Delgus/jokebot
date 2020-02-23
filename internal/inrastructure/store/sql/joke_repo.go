package sql

import (
	"fmt"
	"io"

	"github.com/delgus/jokebot/internal/app"
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

func NewJokeRepo(db *goqu.Database) *JokeRepo {
	return &JokeRepo{db: db}
}

func (j *JokeRepo) GetNewJoke(userID int) (*app.Joke, error) {
	var row jokeRow
	found, err := j.db.From("joke").
		Select("joke.id", "text").
		LeftJoin(goqu.T("joke_user"),
			goqu.On(
				goqu.Ex{"joke_user.joke_id": goqu.I("joke.id")},
				goqu.I("user_id").Eq(userID))).
		Where(goqu.I("user_id").IsNull()).
		ScanStruct(&row)
	if err != nil {
		return nil, fmt.Errorf("infrastructure.store.sql: can't get joke: %v", err)
	}
	if !found {
		return nil, io.EOF
	}
	insert := j.db.Insert("joke_user").Rows(
		goqu.Record{"user_id": userID, "joke_id": row.ID},
	).Executor()

	if _, err := insert.Exec(); err != nil {
		return nil, fmt.Errorf("infrastructure.store.sql: can't save user joke link: %v", err)
	}
	return resolveJokeRow(row), nil
}

func (j *JokeRepo) GetNewJokeByCategory(userID int, categoryID int) (*app.Joke, error) {
	var row jokeRow
	found, err := j.db.From("joke").
		Select("joke.id", "text").
		LeftJoin(goqu.T("joke_user"),
			goqu.On(
				goqu.Ex{"joke_user.joke_id": goqu.I("joke.id")},
				goqu.I("user_id").Eq(userID))).
		Where(goqu.I("user_id").IsNull(), goqu.I("category_id").Eq(categoryID)).
		ScanStruct(&row)
	if err != nil {
		return nil, fmt.Errorf("infrastructure.store.sql: can't get joke: %v", err)
	}
	if !found {
		return nil, io.EOF
	}
	insert := j.db.Insert("joke_user").Rows(
		goqu.Record{"user_id": userID, "joke_id": row.ID},
	).Executor()

	if _, err := insert.Exec(); err != nil {
		return nil, fmt.Errorf("infrastructure.store.sql: can't save user joke link: %v", err)
	}
	return resolveJokeRow(row), nil
}

func (j *JokeRepo) GetJokeCategoryList() ([]app.JokeCategory, error) {
	var rows []jokeCategoryRow
	err := j.db.From("joke_category").
		Select(&jokeCategoryRow{}).
		ScanStructs(&rows)
	if err != nil {
		return nil, fmt.Errorf("infrastructure.store.sql: can't get joke category list: %v", err)
	}
	return resolveJokeCategoryList(rows), nil
}

func resolveJokeRow(row jokeRow) *app.Joke {
	joke := new(app.Joke)
	joke.ID = row.ID
	joke.Text = row.Text
	return joke
}

func resolveJokeCategoryList(rows []jokeCategoryRow) []app.JokeCategory {
	list := make([]app.JokeCategory, len(rows))
	for i, r := range rows {
		list[i] = app.JokeCategory{
			ID:   r.ID,
			Name: r.Name,
		}
	}
	return list
}
