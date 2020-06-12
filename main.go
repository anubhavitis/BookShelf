package main

import (
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseFiles("index.html"))

//Book structure
type Book struct {
	Name    string
	Author  string
	Content string
	Likes   int
}

//IndexHandler function
func IndexHandler(w http.ResponseWriter, req *http.Request) {
	tpl.Execute(w, nil)
}

func main() {
	mux := http.NewServeMux()
	assets := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/submit", IndexHandler)

	// port := os.Getenv("PORT")
	// http.ListenAndServe(":"+port, mux)
	http.ListenAndServe(":8080", mux)
}
