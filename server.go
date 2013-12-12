package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"fmt"
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
	Password []byte
	Email    string
}

func newUser(username, password, email string) User {
	u := User{
		Username: username,
		Email:    email,
	}
	u.SetPassword(password)
	return u
}

// SetPassword takes a plaintext password, hashes it with bcrypt and sets the
// password field to the hash.
func (u *User) SetPassword(password string) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err) // this is a panic because bcrypt errors on invalid costs
	}
	u.Password = hpass
}

// Login validates and returns a user object if they exist in the database.
func Login(dbmap *gorp.DbMap, username, password string) (u *User, err error) {
	err = dbmap.SelectOne(&u, "select * from posts where username=?", username)
	fmt.Println(u)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		u = nil
	}
	return
}

func main() {
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	// create two users
	u1 := newUser("bg", "password", "benzguo@gmail.com")
	u2 := newUser("bzg", "Password2", "ben@venmo.com")

	// insert rows - auto increment PKs will be set properly
	err = dbmap.Insert(&u1, &u2)
	checkErr(err, "Insert failed")

	db := &dbmap
	m := martini.Classic()
	m.Map(db)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/login", func() string {
		u, err := Login(dbmap, "bg", "password")
		fmt.Println(u)
		fmt.Println(err)
		return "Login"
	})
	m.Run()
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
