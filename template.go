package main

import (
	"html/template"
	"net/http"
)

type TemplateData struct {
	ActivePage string
	Flash      *Flash
	Data       interface{}
}

func renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	filename string,
	activePage string,
	data interface{},
) {
	tmpl, err := template.ParseFiles(
		"templates/"+filename,
		"templates/partials/navbar.html",
		"templates/partials/footer.html",
		"templates/partials/flash.html",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := TemplateData{
		ActivePage: activePage,
		Flash:      GetFlash(r, w),
		Data:       data,
	}

	err = tmpl.Execute(w, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
