package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./data"
)

func profile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for profile...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/profile.html"))

	sess, err := session(r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		User data.User
	}

	infos := TempStruct{user}
	templates.ExecuteTemplate(w, "layout", infos)
}
