package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"alkalax/skip-shutdown-html/internal/handler"
	"alkalax/skip-shutdown-html/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		ex, err := os.Executable()
		if err != nil {
			log.Fatalf("failed to get executable path: %v", err)
		}
		templateDir = filepath.Join(filepath.Dir(ex), "templates")
	}

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		log.Fatalf("template directory %s does not exist", templateDir)
	}

	tpl := template.Must(template.ParseFiles(
		filepath.Join(templateDir, "base.html"),
		filepath.Join(templateDir, "partials/dropdown.html"),
	))

	vmHandler := handler.NewVMHandler(tpl, storage.NewMockVMRepository())

	http.HandleFunc("/", vmHandler.GetVM)
	http.HandleFunc("/skip", vmHandler.PostSkipShutdown)

	log.Printf("Server is running on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
