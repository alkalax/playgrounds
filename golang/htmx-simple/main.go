package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title   string
	Heading string
	Message string
}

func main() {
	tpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Alkalax | Home",
			Heading: "Page heading",
			Message: "This is paragraph text.",
		}

		tpl.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}
