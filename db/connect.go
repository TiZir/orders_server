package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("postgres", os.Getenv("PG_URL"))
}

// const (
// 	host     = "postgres"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "postgres"
// )

// func GetDB() (*sql.DB, error) {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Printf("err 1 %v", err)
// 		return db, err
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		log.Printf("err 2 %v", err)
// 		return db, err
// 	}
// 	//return sql.Open("postgres", psqlInfo)
// 	return db, nil
// }
