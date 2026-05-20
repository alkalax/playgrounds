package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		w.Write([]byte("Hello, World!"))
	})

	log.Printf("Server is running on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
