package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("PG_URL"))
}
