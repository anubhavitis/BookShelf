package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/anubhavitis/BookShelf/auth"
	"github.com/anubhavitis/BookShelf/database"
	_ "github.com/go-sql-driver/mysql"
)

var tpl = template.Must(template.ParseFiles("assets/welcome.html"))
var tplauth = template.Must(template.ParseFiles("assets/index.html"))
var db *sql.DB

//IndexHandler function
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login Screen.")
	// books := database.ReadBooks(db)
	// cval, err := auth.ReadCookie(req)
	// if err != nil {
	// 	fmt.Println("Error while reading Cookie")
	// 	return
	// }
	// fmt.Println("the Cookie value:", cval)

	// if auth.CheckSession(cval["sessionID"], req) == true {
	// 	fmt.Println("User auto LoggedIn")
	// } else {
	// 	fmt.Println("No user data found!")
	// 	tplauth.Execute(w, nil)
	// }
	if e := tplauth.Execute(w, nil); e != nil {
		fmt.Println("Template not executed")
		log.Fatalln(e)
	}
}

//SignoutHandler function
func SignoutHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login Page!")
	tplauth.Execute(w, nil)
}

//SubmitHandler function
func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	books := database.ReadBooks(db)
	fmt.Println("Current database ", len(books), " books")

	newBook := &database.Book{}
	newBook.Name = req.FormValue("name")
	newBook.Author = req.FormValue("author")
	newBook.Content = req.FormValue("content")
	newBook.Favo = 0
	fmt.Println("Book read", newBook.Name)

	//If Book is already present
	if database.PriBooks[newBook.Name].Name == newBook.Name {
		fmt.Println("Book already exists in the record.")
		if err := tpl.Execute(w, books); err != nil {
			fmt.Println("Error here at executing template.")
			log.Fatal(err)
		}
		return
	}
	database.AddNewBook(db, *newBook)
	books = append(books, *newBook)

	if e := tpl.Execute(w, books); e != nil {
		fmt.Println("Error here at executing template.")
		log.Fatal(e)
	}
}

//UpdateHandler function
func UpdateHandler(w http.ResponseWriter, req *http.Request) {

	u, err := url.Parse(req.URL.String())
	if err != nil {
		log.Fatal(err)
		return
	}
	params := u.Query()
	name := params.Get("name")
	fav, err1 := strconv.Atoi(params.Get("fav"))
	if err1 != nil {
		log.Fatal(err1)
	}
	del, err2 := strconv.Atoi(params.Get("del"))
	if err2 != nil {
		log.Fatal(err1)
	}
	if del == 0 {
		fav = fav ^ 1
		q := "UPDATE mybooks SET favo=? WHERE name=?"
		if _, e := db.Exec(q, fav, name); e != nil {
			fmt.Println("Upating database error!")
			log.Fatal(e)
		}
	} else {
		q := "DELETE FROM mybooks WHERE name=?"
		if _, e := db.Exec(q, name); e != nil {
			fmt.Println("Upating database error!")
			log.Fatal(e)
		}
		delete(database.PriBooks, name)
	}
	books := database.ReadBooks(db)

	if e := tpl.Execute(w, books); e != nil {
		fmt.Println("Error here at executing template.")
		log.Fatal(e)
	}

}

//SignIn func to for new session auth
func SignIn(w http.ResponseWriter, req *http.Request) {
	tpl.Execute(w, nil)
}

//SignUp func to handle new registration.
func SignUp(w http.ResponseWriter, req *http.Request) {

	newMem := &database.Member{}
	newMem.Name = req.FormValue("name")
	newMem.Email = req.FormValue("email")
	newMem.Password = req.FormValue("password")
	newMem.UID = database.AddMember(db, *newMem)
	fmt.Println("Registering", newMem.Name, "at", newMem.UID)

	if e := auth.CreateCookie(newMem.UID, database.generateUUID, w); e != nil {
		log.Fatalln(e)
	}
	if e := auth.CreateSession(newMem.UID, database.generateUUID, w, req); e != nil {
		log.Fatalln(e)
	}
	tpl.Execute(w, nil)
}

func main() {
	db = database.InitDb()
	database.NewTable(db)

	mux := http.NewServeMux()
	assets := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/submit", SubmitHandler)
	mux.HandleFunc("/update", UpdateHandler)
	mux.HandleFunc("/signout", SignoutHandler)
	mux.HandleFunc("/signin", SignIn)
	mux.HandleFunc("/signup", SignUp)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, mux)

}
