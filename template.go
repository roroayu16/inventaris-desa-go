package main

import (
	"html/template"
	"net/http"
)

type TemplateData struct {
	ActivePage string
	Data       interface{}
}

func renderTemplate(
	w http.ResponseWriter,
	filename string,
	activePage string,
	data interface{},
) {
	tmpl, err := template.ParseFiles(
		"templates/"+filename,
		"templates/partials/navbar.html",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := TemplateData{
		ActivePage: activePage,
		Data:       data,
	}

	err = tmpl.Execute(w, page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
