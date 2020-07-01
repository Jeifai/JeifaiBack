package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func results(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for results...")

	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/topbar.html",
			"templates/sidebar.html",
			"templates/results.html"))

	sess, err := session(r)
	user, err := UserById(sess.UserId)
	if err != nil {
		danger(err, "Cannot find user")
	}
	results, err := user.ResultsByUser()

	type TempStruct struct {
		User User
		Data []Result
	}

	infos := TempStruct{user, results}
	templates.ExecuteTemplate(w, "layout", infos)
}
