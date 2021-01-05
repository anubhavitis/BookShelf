package main

import (
	"database/sql"
	"fmt"
	"html/template"
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
		fmt.Println("1. ", err)
		auth.DeleteCookie(w)
		tplauth.Execute(w, nil)
		return
	}
	f := 0
	if cval == nil {
		f = 1
	} else {
		if a, e := auth.CheckSession(cval["sessionID"], req); e == nil {
			if a == false {
				f = 1
			}
		} else {
			fmt.Println("2 ", e)
			return
		}
	}

	if f == 1 {
		auth.DeleteCookie(w)
		tplauth.Execute(w, nil)
		return
	}
	user, err := database.GetUser(db, cval["userID"])
	if err != nil {
		fmt.Println(err)
		return
	}
	books, err := database.ReadBooks(db, user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user, len(books))
	tpl.Execute(w, books)
}

//SubmitHandler function
func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	cval, err := auth.ReadCookie(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	user, err := database.GetUser(db, cval["userID"])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user, "holding the session on the browser.")
	books, err := database.ReadBooks(db, user)
	if err != nil {
		fmt.Println(err)
		return
	}
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
			fmt.Print(err)
			return
		}
		return
	}
	err1 := database.AddNewBook(db, *newBook, user)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	books = append(books, *newBook)
	// database.ReadAllBooks(db)
	tpl.Execute(w, books)
}

//UpdateHandler function
func UpdateHandler(w http.ResponseWriter, req *http.Request) {
	cval, err := auth.ReadCookie(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := 0
	if cval == nil {
		f = 1
	} else {
		if a, e := auth.CheckSession(cval["sessionID"], req); e == nil {
			if a == false {
				f = 1
			}
		} else {
			fmt.Println(e)
			return
		}
	}
	if f == 1 {
		auth.DeleteCookie(w)
		tplauth.Execute(w, nil)
		return
	}

	u, err := url.Parse(req.URL.String())
	if err != nil {
		fmt.Print(err)
		return
	}

	params := u.Query()
	bookname := params.Get("name")
	fav, err1 := strconv.Atoi(params.Get("fav"))
	if err1 != nil {
		fmt.Print(err1)
		return
	}
	del, err2 := strconv.Atoi(params.Get("del"))
	if err2 != nil {
		fmt.Print(err1)
		return
	}

	cval, er := auth.ReadCookie(req)
	if er != nil {
		fmt.Println(er)
		return
	}
	user, err := database.GetUser(db, cval["userID"])
	if err != nil {
		fmt.Println(err)
		return
	}

	if del == 0 {
		fav = fav ^ 1
		q := "UPDATE mybooks SET favo=? WHERE name=? and user=?"
		if _, e := db.Exec(q, fav, bookname, user); e != nil {
			fmt.Println("Upating database error!")
			fmt.Print(e)
			return
		}
	} else {
		q := "DELETE FROM mybooks WHERE name=? and user=?"
		if _, e := db.Exec(q, bookname, user); e != nil {
			fmt.Println("Upating database error!")
			fmt.Print(e)
			return
		}
		delete(database.PriBooks, bookname)
	}
	books, err := database.ReadBooks(db, user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(books))

	tpl.Execute(w, books)

}

//SignIn func to for new session auth
func SignIn(w http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	password := req.FormValue("password")

	userID, authkey, err := database.GetIDPassword(db, email)
	if err != nil {
		fmt.Println(err)
		return
	}

	if authkey != password {
		fmt.Println("Wrong Password")
		return
	}
	fmt.Println("User Authenticated!")
	sID := database.GenerateUUID()
	if e := auth.CreateCookie(userID, sID, w); e != nil {
		fmt.Println("Error while creating cookie")
		fmt.Println(e)
		return
	}
	if e := auth.CreateSession(userID, sID, w, req); e != nil {
		fmt.Println("Error while creating a session")
		fmt.Println(e)
		return
	}
	user, err := database.GetUser(db, userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	books, err := database.ReadBooks(db, user)
	if err != nil {
		fmt.Println(err)
		return
	}
	tpl.Execute(w, books)
}

//SignUp func to handle new registration.
func SignUp(w http.ResponseWriter, req *http.Request) {

	newMem := &database.Member{}
	newMem.Name = req.FormValue("name")
	newMem.Email = req.FormValue("email")
	newMem.Password = req.FormValue("password")

	name, err := database.AddMember(db, *newMem)
	if err != nil {
		fmt.Println(err)
		return
	}
	newMem.UID = name
	fmt.Println("Registering", newMem.Name, "at", newMem.UID)
	sID := database.GenerateUUID()
	if e := auth.CreateCookie(newMem.UID, sID, w); e != nil {
		fmt.Println("Error while creating cookie")
		fmt.Println(e)
		return
	}
	if e := auth.CreateSession(newMem.UID, sID, w, req); e != nil {
		fmt.Println("Error while creating a session")
		fmt.Println(e)
		return
	}
	user, err := database.GetUser(db, newMem.UID)
	if err != nil {
		fmt.Println(err)
		return
	}
	books, err := database.ReadBooks(db, user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user, "holding the session on the browser.")
	tpl.Execute(w, books)
}

//SignoutHandler function
func SignoutHandler(w http.ResponseWriter, req *http.Request) {
	auth.DeleteCookie(w)
	auth.ClearSession(w, req)
	tplauth.Execute(w, nil)
}

func main() {
	dab, err := database.InitDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	db = dab
	database.NewBookTable(db)
	database.NewMemberTable(db)

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
