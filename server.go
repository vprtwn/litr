package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/coopernurse/gorp"
	"github.com/lib/pq"
	"log"
	"net/http"
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

func NewUser(username, password, email string) User {
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
	checkErr(err, "bcrypt failed")
	u.Password = hpass
}

// LogIn validates and returns a user object if they exist in the database.
func LogIn(dbmap *gorp.DbMap, username, password string) (u *User, err error) {
	var us []User
	_, err = dbmap.Select(&us, "select * from users where Username = :username",
		map[string]interface{}{"username": username})
	checkErr(err, "Select users with matching Username failed")
	if len(us) == 0 {
		err = errors.New("No user with matching username found")
	} else if len(us) > 1 {
		err = errors.New("Multiple users with the same username found")
	} else {
		u = &us[0]
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		u = nil
	}
	return
}

/* u, err := Login(dbmap, "bg", "password") */
/* log.Println("user: %v", u) */
/* checkErr(err, "Login error") */
/* return "Login" */

func main() {
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	// create two users
	u1 := NewUser("bg", "password", "benzguo@gmail.com")
	u2 := NewUser("bzg", "Password2", "ben@venmo.com")

	// insert rows - auto increment PKs will be set properly
	err = dbmap.Insert(&u1, &u2)
	checkErr(err, "Insert failed")

	db := &dbmap
	m := martini.Classic()
	m.Map(db)
	// render html templates from templates directory
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "home", nil)
	})

	m.Post("/login", func(w http.ResponseWriter, req *http.Request, r render.Render) {
		u, err := LogIn(dbmap, req.FormValue("username"), req.FormValue("password"))
		if err != nil {
			//TODO: session flash with gorilla
			log.Println(err)
			r.HTML(200, "home", nil)
		}
		if u != nil {
			http.Redirect(w, req, "/u/"+u.Username, http.StatusFound)
		}
	})

	m.Get("/u/:username", func(params martini.Params) string {
		return "Hello " + params["username"]
	})

	m.Run()
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
