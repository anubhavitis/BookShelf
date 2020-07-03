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
	
    // auth package import is commented because auth is not used anywhere in the file
	// "github.com/anubhavitis/BookShelf/auth" 
	"github.com/anubhavitis/BookShelf/database"

	_ "github.com/go-sql-driver/mysql"
)

var tpl = template.Must(template.ParseFiles("welcome.html"))
var tpl0 = template.Must(template.ParseFiles("index.html"))
var db *sql.DB
var priBooks = make(map[string]Book)
var priMember = make(map[string]Member)

//Book structure
type Book struct {
	Name    string
	Author  string
	Content string
	Favo    int
}

//Member stuct for registered users
type Member struct {
	Name     string
	Email    string
	Password string
}

//ReadBooks read all books and stores it in the map books
func ReadBooks() []Book {
	var books []Book
	rows, er := db.Query(`SELECT name,author,content,favo FROM mybooks`)
	if er != nil {
		log.Fatal(er)
	}
	defer rows.Close()

	for rows.Next() {
		var temp Book
		if err := rows.Scan(&temp.Name, &temp.Author, &temp.Content, &temp.Favo); err != nil {
			log.Fatal(err)
		}
		books = append(books, temp)
		priBooks[temp.Name] = temp
	}
	return books
}

//IndexHandler function
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Home is reached.")
	books := ReadBooks()
	tpl.Execute(w, books)
}

//LoginHandler function
func LoginHandler(w http.ResponseWriter, req *http.Request) {
	tpl0.Execute(w, nil)
}

//SubmitHandler function
func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	u, err := url.Parse(req.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	books := ReadBooks()
	fmt.Println("Current database")
	for x := range books {
		fmt.Println(books[x].Name, books[x].Author, books[x].Content, books[x].Favo)
	}
	params := u.Query()
	newBook := &Book{}
	newBook.Name = params.Get("name")
	newBook.Author = params.Get("author")
	newBook.Content = params.Get("content")
	newBook.Favo = 0
	fmt.Println("Book read", newBook.Name)

	//If Book is already present
	if priBooks[newBook.Name].Name == newBook.Name {
		fmt.Println("Book already exists in the record.")
		if err := tpl.Execute(w, books); err != nil {
			fmt.Println("Error here at executing template.")
			log.Fatal(err)
		}
		return
	}

	query := ` INSERT INTO mybooks (name, author, content, favo) VALUES (?,?,?,?,?)`
	if _, e := db.Exec(query, newBook.Name, newBook.Author, newBook.Content, newBook.Favo); e != nil {
		log.Fatal(err)
	}
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
		delete(priBooks, name)
	}
	books := ReadBooks()

	if e := tpl.Execute(w, books); e != nil {
		fmt.Println("Error here at executing template.")
		log.Fatal(e)
	}

}

// //SignIn func to for new session auth
// func SignIn(w http.ResponseWriter, req *http.ResponseWriter) {
// 	u,e:=url.Parse(req.URL.String())
// 	if()
// }

//SignUp func to handle new registration.
func SignUp(w http.ResponseWriter, req *http.Request) {
	u, e := url.Parse(req.URL.String())
	if e != nil {
		log.Fatal(e)
		return
	}
	params := u.Query()
	var newM Member
	newM.Name = params.Get("name")
	newM.Email = params.Get("email")
	newM.Password = params.Get("password")
	if priMember[newM.Email].Email == newM.Email {
		fmt.Println("User Already Exists")
		if e := tpl0.Execute(w, nil); e != nil {
			log.Fatal(e)
		}
		return
	}

	q := ` INSERT INTO members (name, email, password) VALUES (?,?,?)`
	if _, e := db.Exec(q, newM.Name, newM.Email, newM.Password); e != nil {
		fmt.Println("Upating database error!")
		log.Fatal(e)
	}
	priMember[newM.Email] = newM
	books := ReadBooks()
	if e := tpl.Execute(w, books); e != nil {
		log.Fatal(e)
	}
}

func main() {
	db = database.InitDb()
	database.NewTable(db)

	mux := http.NewServeMux()
	assets := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/submit", SubmitHandler)
	mux.HandleFunc("/update", UpdateHandler)
	mux.HandleFunc("/Welcome", LoginHandler)
	// mux.HandleFunc("/signin", SignIn)
	mux.HandleFunc("/signup", SignUp)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, mux)

}
