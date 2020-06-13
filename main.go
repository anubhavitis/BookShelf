package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

var tpl = template.Must(template.ParseFiles("index.html"))
var db *sql.DB
var priBooks = make(map[string]bool)

func initDb() {
	dab, err := sql.Open("mysql", "zeddie:1234567890@(127.0.0.1:3306)/mysql?parseTime=true")
	if err != nil {
		fmt.Println("Error at opening database")
		log.Fatal(err)
	}
	if err := dab.Ping(); err != nil {
		fmt.Println("Error at ping.")
		log.Fatal(err)
	}
	db = dab
}

func newTable() {

	if _, err := db.Exec("DROP TABLE mybooks"); err != nil {
		log.Fatal(err)
	}

	query := `
	CREATE TABLE mybooks(
		id INT AUTO_INCREMENT,
		name TEXT NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		likes INT,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
	fmt.Println("New Table Created!")
}

//Book structure
type Book struct {
	Name    string
	Author  string
	Content string
	Likes   int
}

//ReadBooks read all books and stores it in the map books
func ReadBooks() []Book {
	var books []Book
	rows, er := db.Query(`SELECT name,author,content,likes FROM mybooks`)
	if er != nil {
		log.Fatal(er)
	}
	defer rows.Close()

	for rows.Next() {
		var temp Book
		if err := rows.Scan(&temp.Name, &temp.Author, &temp.Content, &temp.Likes); err != nil {
			log.Fatal(err)
		}
		books = append(books, temp)
		priBooks[temp.Name] = true
	}
	return books
}

//IndexHandler function
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Home is reached.")
	books := ReadBooks()
	tpl.Execute(w, books)
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
		fmt.Println(books[x].Name, books[x].Author, books[x].Content, books[x].Likes)
	}
	params := u.Query()
	newBook := &Book{}
	newBook.Name = params.Get("name")
	newBook.Author = params.Get("author")
	newBook.Content = params.Get("content")
	newBook.Likes = 0
	fmt.Println("Book read", newBook.Name)

	//If Book is already present
	if priBooks[newBook.Name] == true {
		fmt.Println("Book already exists in the record.")
		if err := tpl.Execute(w, books); err != nil {
			fmt.Println("Error here at executing template.")
			log.Fatal(err)
		}
		return
	}

	query := ` INSERT INTO mybooks (name, author, content, likes) VALUES (?,?,?,?)`
	if _, e := db.Exec(query, newBook.Name, newBook.Author, newBook.Content, newBook.Likes); e != nil {
		log.Fatal(err)
	}
	books = append(books, *newBook)

	if e := tpl.Execute(w, books); e != nil {
		fmt.Println("Error here at executing template.")
		log.Fatal(e)
	}
}

func main() {
	initDb()
	newTable()

	mux := http.NewServeMux()
	assets := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/submit", SubmitHandler)

	// port := os.Getenv("PORT")
	// http.ListenAndServe(":"+port, mux)
	http.ListenAndServe(":8080", mux)
}
