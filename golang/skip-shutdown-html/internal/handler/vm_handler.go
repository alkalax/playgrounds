package handler

import (
	"alkalax/skip-shutdown-html/internal/models"
	"alkalax/skip-shutdown-html/internal/service"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)

var vms = []models.VirtualMachineInfo{
	{Name: "VM1", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
	{Name: "VM2", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: true},
	{Name: "VM3", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
}

type VMHandler struct {
	tpl *template.Template
}

func NewVMHandler(tpl *template.Template) *VMHandler {
	return &VMHandler{tpl: tpl}
}

func (h *VMHandler) GetVM(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	selectedName := r.URL.Query().Get("vm")
	pageData := models.PageData{
		VMs:          vms,
		SelectedName: selectedName,
		Selected:     service.FindVM(selectedName),
	}

	h.tpl.ExecuteTemplate(w, "base", pageData)
}

func (h *VMHandler) PostSkipShutdown(w http.ResponseWriter, r *http.Request) {
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
	if vm := service.FindVM(selectedName); vm != nil {
		vm.SkipToday = true
	}

	http.Redirect(w, r, "/?vm="+url.QueryEscape(selectedName), http.StatusSeeOther)
}
