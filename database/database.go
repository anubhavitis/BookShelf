package database

import (
	"database/sql"
	"fmt"
	"log"
)

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

	query = `
	CREATE TABLE members(
		id INT AUTO_INCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
	);`
	if _, err := db.Exec(query); err != nil {
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
