package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(c string) error {

	var e error
	db, e = sql.Open("postgres", c)
	if e != nil {
		return e
	}

	return nil
}
