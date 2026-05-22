package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"alkalax/skip-shutdown-html/internal/models"
)

var vms = []models.VirtualMachineInfo{
	{Name: "VM1", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
	{Name: "VM2", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: true},
	{Name: "VM3", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
}

func findVM(name string) *models.VirtualMachineInfo {
	for i := range vms {
		if vms[i].Name == name {
			return &vms[i]
		}
	}
	return nil
}

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		selectedName := r.URL.Query().Get("vm")
		pageData := models.PageData{
			VMs:          vms,
			SelectedName: selectedName,
			Selected:     findVM(selectedName),
		}

		tpl.ExecuteTemplate(w, "base", pageData)
	})

	http.HandleFunc("/skip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("failed to parse form: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		selectedName := r.FormValue("vm")
		if vm := findVM(selectedName); vm != nil {
			vm.SkipToday = true
		}

		http.Redirect(w, r, "/?vm="+url.QueryEscape(selectedName), http.StatusSeeOther)
	})

	log.Printf("Server is running on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
