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

	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/keywords.html"))

	type TempStruct struct {
		User    data.User
		Targets []string
	}

	struct_targets, err := user.UsersTargetsByUser()

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Url)
	}

	infos := TempStruct{user, arr_targets}
	templates.ExecuteTemplate(w, "layout", infos)
}
