package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./data"
)

func keywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for keywords...")
	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		danger(err, "Cannot find user")
    }
    _ = user
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/keywords.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}
