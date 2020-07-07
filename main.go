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

	cval, err := auth.ReadCookie(req)
	if err != nil {
		log.Fatalln(err)
	}
	f := 0
	if cval["userID"] == "" {
		f = 1
	} else if auth.CheckSession(cval["sessionID"], req) == false {
		f = 1
	}

	if f == 1 {
		tplauth.Execute(w, nil)
		return
	}

	books := database.ReadBooks(db, database.GetUser(db, cval["userID"]))
	tpl.Execute(w, books)
}

//SubmitHandler function
func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	cval, err := auth.ReadCookie(req)
	if err != nil {
		log.Fatalln(err)
	}
	user := database.GetUser(db, cval["userID"])
	books := database.ReadBooks(db, user)
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
	database.AddNewBook(db, *newBook, cval["userID"])
	books = append(books, *newBook)

	tpl.Execute(w, books)
}

//UpdateHandler function
func UpdateHandler(w http.ResponseWriter, req *http.Request) {

	u, err := url.Parse(req.URL.String())
	if err != nil {
		log.Fatal(err)
		return
	}
	params := u.Query()
	bookname := params.Get("name")
	fav, err1 := strconv.Atoi(params.Get("fav"))
	if err1 != nil {
		log.Fatal(err1)
	}
	del, err2 := strconv.Atoi(params.Get("del"))
	if err2 != nil {
		log.Fatal(err1)
	}
	cval, err := auth.ReadCookie(req)
	if err != nil {
		log.Fatalln(err)
	}
	user := database.GetUser(db, cval["userID"])
	if del == 0 {
		fav = fav ^ 1
		q := "UPDATE mybooks SET favo=? WHERE name=? and user=?"
		if _, e := db.Exec(q, fav, bookname, user); e != nil {
			fmt.Println("Upating database error!")
			log.Fatal(e)
		}
	} else {
		q := "DELETE FROM mybooks WHERE name=? and user=?"
		if _, e := db.Exec(q, bookname, user); e != nil {
			fmt.Println("Upating database error!")
			log.Fatal(e)
		}
		delete(database.PriBooks, bookname)
	}
	books := database.ReadBooks(db, user)
	tpl.Execute(w, books)

}

//SignIn func to for new session auth
func SignIn(w http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	password := req.FormValue("password")

	userID, authkey := database.GetIDPassword(db, email)

	if authkey != password {
		fmt.Println("Wrong Password")
		return
	}
	fmt.Println("User Authenticated!")
	sID := database.GenerateUUID()
	if e := auth.CreateCookie(userID, sID, w); e != nil {
		fmt.Println("Error while creating cookie")
		log.Fatalln(e)
	}
	if e := auth.CreateSession(userID, sID, w, req); e != nil {
		fmt.Println("Error while creating a session")
		log.Fatalln(e)
	}
	books := database.ReadBooks(db, database.GetUser(db, userID))
	tpl.Execute(w, books)
}

//SignUp func to handle new registration.
func SignUp(w http.ResponseWriter, req *http.Request) {

	newMem := &database.Member{}
	newMem.Name = req.FormValue("name")
	newMem.Email = req.FormValue("email")
	newMem.Password = req.FormValue("password")
	newMem.UID = database.AddMember(db, *newMem)
	fmt.Println("Registering", newMem.Name, "at", newMem.UID)
	sID := database.GenerateUUID()
	if e := auth.CreateCookie(newMem.UID, sID, w); e != nil {
		fmt.Println("Error while creating cookie")
		log.Fatalln(e)
	}
	if e := auth.CreateSession(newMem.UID, sID, w, req); e != nil {
		fmt.Println("Error while creating a session")
		log.Fatalln(e)
	}
	books := database.ReadBooks(db, database.GetUser(db, newMem.UID))
	tpl.Execute(w, books)
}

//SignoutHandler function
func SignoutHandler(w http.ResponseWriter, req *http.Request) {
	auth.DeleteCookie(w)
	auth.ClearSession(w, req)
	tplauth.Execute(w, nil)
}

func main() {
	db = database.InitDb()
	// database.NewBookTable(db)
	// database.NewMemberTable(db)

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
