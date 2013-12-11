package main

import (
	"database/sql"
	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	"github.com/lib/pq"
	"log"
	"os"
)

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	connection += " sslmode=disable"

	db, err := sql.Open("postgres", connection)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'users' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

type User struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id       int64 `db:"user_id"`
	Username string
	Password string
	Email    string
}

func newUser(username, password, email string) User {
	return User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

func main() {
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	db := &dbmap
	m := martini.Classic()
	m.Map(db)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
