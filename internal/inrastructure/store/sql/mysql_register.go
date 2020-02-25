package sql

import (
	// register mysql addr
	_ "github.com/doug-martin/goqu/v8/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)
