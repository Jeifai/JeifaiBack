package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func jobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for jobs...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/jobs.html"))

	templates.ExecuteTemplate(w, "layout", nil)
}