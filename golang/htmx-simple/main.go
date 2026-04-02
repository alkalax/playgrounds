package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title   string
	Heading string
	Message string
	Items   []string
}

func main() {
	tpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/content.html",
		"templates/partials/message.html",
		"templates/partials/items.html",
	))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Alkalax | Home",
			Heading: "Page heading",
			Message: "This is paragraph text.",
			Items:   []string{"apples", "oranges", "pears", "watermelons"},
		}

		tpl.ExecuteTemplate(w, "base", data)
	})

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Items: []string{"raspberries", "blueberries", "pears"},
		}

		tpl.ExecuteTemplate(w, "items", data)
	})

	http.ListenAndServe(":8080", nil)
}
