package main

import (
	"./data"
	"fmt"
	"html/template"
	"net/http"
)

func results(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for results...")
	sess, err := session(r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
	results, err := user.ResultsByUser()
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/results.html"))
	type TempStruct struct {
		User    data.User
		Results []data.Result
		Message string
	}
	infos := TempStruct{user, results, "Here the list of all your results"}
	templates.ExecuteTemplate(w, "layout", infos)
}
