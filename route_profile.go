package main

import (
    "fmt"
    "net/http"
	"html/template"
)

func profile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Generating HTML for profile...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/profile.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}