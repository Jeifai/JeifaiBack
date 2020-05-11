package main

import (
    "./data"
    "fmt"
	"net/http"
)

func targets(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "layout", "private.navbar", "targets")
}

func target_add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for target...")
	generateHTML(w, nil, "layout", "private.navbar", "target_add")
}

func target_add_url(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Adding URL into target")
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	url := data.Url{
		Url: r.PostFormValue("url"),
	}
	if err := url.CreateTarget(); err != nil {
		danger(err, "Cannot create url")
	}
	http.Redirect(w, r, "/targets", 302)
}