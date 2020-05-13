package main

import (
	"./data"
	"fmt"
	"html/template"
	"net/http"
)

func targets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for targets...")
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
	targets, err := user.UsersTargets()

	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/targets.html"))

	type TempStruct struct {
		User    data.User
		Targets []data.Target
	}
	infos := TempStruct{user, targets}

	templates.ExecuteTemplate(w, "layout", infos)
}

func target_add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for target_add...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/target_add.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func target_add__run(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting target_add_url...")
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	target := data.Target{
		Url: r.PostFormValue("url"),
	}
	if err := target.CreateTarget(); err != nil {
		danger(err, "Cannot create target")
	}
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
	if err := target.CreateUserTarget(user); err != nil {
		danger(err, "Cannot create relation user <--> target")
	}
	fmt.Println("Closing target_add_url...")
	http.Redirect(w, r, "/targets", 302)
}
