package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/sessions"
	"github.com/coopernurse/gorp"
	"github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func initDb() *gorp.DbMap {
	// connect to db using database/sql & pq
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

	// create the table. In a production system you'd generally
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

type ProfilePage struct {
	Title    string
	LoggedIn bool
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

func main() {
	// initialize database
	dbmap := initDb()
	defer dbmap.Db.Close()
	db := &dbmap

	// --- test data
	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	// create two users
	u1 := NewUser("bg", "password", "benzguo@gmail.com")
	u2 := NewUser("bzg", "Password2", "ben@venmo.com")

	// insert rows - auto increment PKs will be set properly
	err = dbmap.Insert(&u1, &u2)
	checkErr(err, "Insert failed")
	// -- end test data

	// set up cookie store
	store := sessions.NewCookieStore([]byte(os.Getenv("KEY")))

	// set up martini
	m := martini.Classic()
	m.Map(db)
	m.Use(sessions.Sessions("litr", store))

	m.Get("/", func(w http.ResponseWriter, session sessions.Session) {
		t, err := template.ParseFiles("home.html")
		checkErr(err, "Failed to parse template")
		t.Execute(w, nil)
	})

	m.Post("/login", func(w http.ResponseWriter, r *http.Request, session sessions.Session) {
		u, err := LogIn(dbmap, r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			session.AddFlash("Invalid Username/Password")
			log.Println(err)
		}
		if u != nil {
			session.Set("user", strconv.FormatInt(u.Id, 10))
			http.Redirect(w, r, "/u/"+u.Username, http.StatusFound)
		}
	})

	m.Post("/logout", func(w http.ResponseWriter, r *http.Request, session sessions.Session) {
		session.Delete("user")
		http.Redirect(w, r, "/", http.StatusFound)
	})

	m.Get("/u/:username", func(w http.ResponseWriter, params martini.Params, session sessions.Session) {
		id := session.Get("user")
		var p *ProfilePage
		if id == nil {
			p = &ProfilePage{Title: "Not logged in", LoggedIn: false}
		} else {
			p = &ProfilePage{Title: id.(string), LoggedIn: true}
		}
		t, err := template.ParseFiles("user.html")
		checkErr(err, "Failed to parse template")
		t.Execute(w, p)
	})

	m.Run()
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
