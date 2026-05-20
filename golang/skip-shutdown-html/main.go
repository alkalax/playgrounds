package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type VirtualMachineInfo struct {
	Name         string
	ShutdownTime time.Time
	SkipToday    bool
}

type PageData struct {
	VMs          []VirtualMachineInfo
	SelectedName string
	Selected     *VirtualMachineInfo
}

var vms = []VirtualMachineInfo{
	{Name: "VM1", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
	{Name: "VM2", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: true},
	{Name: "VM3", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
}

func findVM(name string) *VirtualMachineInfo {
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

	tpl := template.Must(template.ParseFiles(
		"templates/base.html",
	))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		selectedName := r.URL.Query().Get("vm")
		pageData := PageData{
			VMs:          vms,
			SelectedName: selectedName,
			Selected:     findVM(selectedName),
		}

		tpl.ExecuteTemplate(w, "base", pageData)
	})

	log.Printf("Server is running on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
