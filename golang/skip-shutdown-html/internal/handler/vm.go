package handler

import (
	"alkalax/skip-shutdown-html/internal/models"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type VMRepository interface {
	GetVMs() []models.VirtualMachineInfo
	FindVM(name string) *models.VirtualMachineInfo
}

type VMHandler struct {
	tpl  *template.Template
	repo VMRepository
}

func NewVMHandler(tpl *template.Template, repo VMRepository) *VMHandler {
	return &VMHandler{
		tpl:  tpl,
		repo: repo,
	}
}

func (h *VMHandler) GetVM(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	selectedName := r.URL.Query().Get("vm")
	pageData := models.PageData{
		VMs:          h.repo.GetVMs(),
		SelectedName: selectedName,
		Selected:     h.repo.FindVM(selectedName),
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
	if vm := h.repo.FindVM(selectedName); vm != nil {
		vm.SkipToday = true
	}

	http.Redirect(w, r, "/?vm="+url.QueryEscape(selectedName), http.StatusSeeOther)
}
