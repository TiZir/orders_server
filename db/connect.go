package db

import (
	"database/sql"
	"os"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("PG_URL"))
}
