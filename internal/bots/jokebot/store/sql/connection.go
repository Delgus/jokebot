package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v8"
)

// ConnectionOptions for jokebot sql
type ConnectionOptions struct {
	Driver string
	Addr   string
	Debug  bool
	Logger goqu.Logger
}

// NewConnection return new connection
func NewConnection(o ConnectionOptions) (*goqu.Database, error) {
	con, err := sql.Open(o.Driver, o.Addr)
	if err != nil {
		return nil, fmt.Errorf(`jokebot.store.sql: cannot open db connection: %v`, err)
	}

	err = con.Ping()
	if err != nil {
		return nil, fmt.Errorf("jokebot.store.sql: db ping failed:: %v", err)
	}

	con.SetConnMaxLifetime(time.Hour)

	db := goqu.New(o.Driver, con)
	if o.Debug {
		db.Logger(o.Logger)
	}
	return db, nil
}
