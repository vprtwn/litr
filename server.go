package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func openDb() *sql.DB {
	connection := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
	}

	return db
}

func main() {
	/* db := openDb() */
	m := martini.Classic()
	/* m.Map(db) */
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}
