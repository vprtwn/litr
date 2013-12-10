package main

import (
	"database/sql"
	"fmt"
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
	sqlDB := openDb()
	db := &sqlDB
	m := martini.Classic()
	m.Map(db)
	m.Get("/", func() string {
		fmt.Println(*sqlDB)
		return "Hello world!"
	})
	m.Run()
}
