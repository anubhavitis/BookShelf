package database

import (
	"database/sql"
	"fmt"
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
func InitDb() (*sql.DB, error) {

	dab, err := sql.Open("mysql", "sql12385164:Ij9fwZBM5s@(sql12.freemysqlhosting.net)/sql12385164?parseTime=true")
	if err != nil {
		fmt.Println("Error at opening database")
		return nil, err
	}
	if err := dab.Ping(); err != nil {
		fmt.Println("Error at ping.")
		return nil, err
	}
	return dab, nil
}

//NewBookTable ..
func NewBookTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS
	mybooks(
		id INT AUTO_INCREMENT,
		user TEXT NOT NULL,
		name TEXT NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		favo INT,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	fmt.Println("mybooks Created!")
	return nil
}

//NewMemberTable ..
func NewMemberTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS
	members(
		id INT AUTO_INCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		fmt.Println("Error occured while creating table.")
		return err
	}
	fmt.Println("Members table Created!")
	return nil
}

//ReadBooks read all books and stores it in the map books
func ReadBooks(db *sql.DB, user string) ([]Book, error) {
	var books []Book
	q := `SELECT name,author,content,favo FROM mybooks WHERE user=?`
	rows, er := db.Query(q, user)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	for rows.Next() {
		var temp Book
		if err := rows.Scan(&temp.Name, &temp.Author, &temp.Content, &temp.Favo); err != nil {
			return nil, err
		}
		books = append(books, temp)
		PriBooks[temp.Name] = temp
	}
	return books, nil
}

//ReadAllBooks read all books and stores it in the map books
func ReadAllBooks(db *sql.DB) error {
	rows, er := db.Query(`SELECT id,user,name,author FROM mybooks`)
	if er != nil {
		return er
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, user, author string
		if err := rows.Scan(&id, &user, &name, &author); err != nil {
			return err
		}
		fmt.Println(id, user, name, author)
	}
	return nil
}

//AddMember ..
func AddMember(db *sql.DB, newMem Member) (string, error) {
	q := ` INSERT INTO members
		(name, email,password)
		Values(?,?,?)`

	if _, e := db.Exec(q, newMem.Name, newMem.Email, newMem.Password); e != nil {
		fmt.Println("member not added to record.")
		return "", e
	}
	q = `SELECT id FROM members WHERE email=?`
	rec, e := db.Query(q, newMem.Email)
	if e != nil {
		fmt.Println("Error at query to find a record")
		return "", e
	}
	defer rec.Close()
	var id int
	for rec.Next() {
		if err := rec.Scan(&id); err != nil {
			fmt.Println("Error at scanning resulted query")
			return "", err
		}
	}
	return strconv.Itoa(id), nil
}

//AddNewBook ..
func AddNewBook(db *sql.DB, newBook Book, user string) error {
	query := ` INSERT INTO mybooks 
	(name,user, author, content, favo) 
	VALUES (?,?,?,?,?)`

	if _, e := db.Exec(query, newBook.Name, user, newBook.Author,
		newBook.Content, newBook.Favo); e != nil {
		fmt.Println("Error while adding books.")
		return e
	}
	fmt.Println("Book added to database!")
	return nil
}

//GetIDPassword ..
func GetIDPassword(db *sql.DB, email string) (string, string, error) {
	q := `SELECT id,password FROM members WHERE email=?`
	rec, e := db.Query(q, email)
	if e != nil {
		fmt.Println("Error at query to find a record")
		return "", "", e
	}
	defer rec.Close()

	var pkey string
	var id int
	for rec.Next() {
		if err := rec.Scan(&id, &pkey); err != nil {
			return "", "", err
		}
	}
	return strconv.Itoa(id), pkey, nil
}

//GetUser ..
func GetUser(db *sql.DB, uID string) (string, error) {
	q := `SELECT name FROM members WHERE id=?`
	row, err := db.Query(q, uID)
	if err != nil {
		fmt.Println("Error at query to find a user")
		return "", err
	}
	defer row.Close()
	var username string
	for row.Next() {
		if e := row.Scan(&username); e != nil {
			return "", e
		}
	}
	return username, nil
}
