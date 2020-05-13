package main

import (
	"./data"
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting index...")
	sess, err := session(w, r)
	if err != nil {
		fmt.Println("Generating HTML for index, user not logged in...")
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/public.navigation.html",
				"templates/index.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		user, err := data.UserByEmail(sess.Email)
		if err != nil {
			danger(err, "Cannot find user")
		}
		fmt.Println("Generating HTML for index, user logged in...")
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/private.navigation.html",
				"templates/index.html"))
		type TempStruct struct {
			User data.User
		}
		infos := TempStruct{user}
		templates.ExecuteTemplate(w, "layout", infos)
	}
}
