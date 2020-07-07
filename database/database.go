package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
)

//Member stuct for registered users
type Member struct {
	UID      string
	Name     string
	Email    string
	Password string
}

//GenerateUUID ..
func GenerateUUID() string {
	v, _ := uuid.NewUUID()
	return v.String()
}

//PriBooks ..
var PriBooks = make(map[string]Book)

//PriMember ..
var PriMember = make(map[string]Member)

//Book structure
type Book struct {
	Name    string
	Author  string
	Content string
	Favo    int
}

//InitDb ..
func InitDb() *sql.DB {

	dab, err := sql.Open("mysql", "sql12349917:VEDK9mPCkq@(sql12.freemysqlhosting.net)/sql12349917?parseTime=true")
	if err != nil {
		fmt.Println("Error at opening database")
		log.Fatal(err)
	}
	if err := dab.Ping(); err != nil {
		fmt.Println("Error at ping.")
		log.Fatal(err)
	}
	return dab
}

//NewTable ..
func NewTable(db *sql.DB) {

	if _, err := db.Exec("DROP TABLE mybooks"); err != nil {
		log.Fatal(err)
	}

	query := `
	CREATE TABLE mybooks(
		id INT AUTO_INCREMENT,
		name TEXT NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		favo INT,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
	fmt.Println("mybooks Created!")

	if _, err := db.Exec("DROP TABLE members"); err != nil {
		log.Fatal(err)
	}

	query = `
	CREATE TABLE members(
		id INT AUTO_INCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		fmt.Println("Error occured while creating table.")
		log.Fatal(err)
	}
	fmt.Println("Members table Created!")
}

//ReadBooks read all books and stores it in the map books
func ReadBooks(db *sql.DB) []Book {
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
		PriBooks[temp.Name] = temp
	}
	return books
}

//AddMember ..
func AddMember(db *sql.DB, newMem Member) string {
	q := ` INSERT INTO members
		(name, email,password)
		Values(?,?,?)`

	if _, e := db.Exec(q, newMem.Name, newMem.Email, newMem.Password); e != nil {
		log.Fatalln(e)
		fmt.Println("member not added to record.")
	}
	q = `SELECT id FROM members WHERE email=?`
	rec, e := db.Query(q, newMem.Email)
	if e != nil {
		log.Fatalln(e)
		fmt.Println("Error at query to find a record")
	}
	defer rec.Close()
	var id int
	for rec.Next() {
		if err := rec.Scan(&id); err != nil {
			log.Fatal(err)
			fmt.Println("Error at scanning resulted query")
		}
	}
	return strconv.Itoa(id)
}

//AddNewBook ..
func AddNewBook(db *sql.DB, newBook Book) {
	query := ` INSERT INTO mybooks 
	(name, author, content, favo) 
	VALUES (?,?,?,?)`

	if _, e := db.Exec(query, newBook.Name, newBook.Author,
		newBook.Content, newBook.Favo); e != nil {
		fmt.Println("Error while adding books.")
		log.Fatal(e)
	}
	fmt.Println("Book added to database!")
}

//GetPassword ..
func GetPassword(db *sql.DB, email string) (string, string) {
	q := `SELECT id,password FROM members WHERE email=?`
	rec, e := db.Query(q, email)
	if e != nil {
		log.Fatalln(e)
		fmt.Println("Error at query to find a record")
	}
	defer rec.Close()

	var pkey string
	var id int
	for rec.Next() {
		if err := rec.Scan(&id, &pkey); err != nil {
			log.Fatal(err)
		}
	}
	return strconv.Itoa(id), pkey
}
