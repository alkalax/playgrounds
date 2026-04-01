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
	tpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/content.html",
		"templates/partials/message.html",
	))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Alkalax | Home",
			Heading: "Page heading",
			Message: "This is paragraph text.",
		}

		tpl.ExecuteTemplate(w, "base", data)
	})

	http.ListenAndServe(":8080", nil)
}
