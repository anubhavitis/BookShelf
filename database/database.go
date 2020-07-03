package database

import (
	"database/sql"
	"fmt"
	"log"
)


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
