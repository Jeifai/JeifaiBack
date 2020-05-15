package main

import (
	"./data"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

func targets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for targets...")
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
	targets, err := user.UsersTargetsByUser()
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/targets.html"))
	type TempStruct struct {
		User    data.User
		Targets []data.Target
		Message string
	}
	infos := TempStruct{user, targets, ""}

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
	fmt.Println("Starting target_add__run...")
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	temp_url := r.PostFormValue("url")
	temp_host, err := url.Parse(temp_url)
	target := data.Target{Url: temp_url, Host: temp_host.Host}
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
	targets, err := user.UsersTargetsByUser()
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/targets.html"))
	type TempStruct struct {
		User    data.User
		Targets []data.Target
		Message string
	}
	infos := TempStruct{user, targets, "Target successfully added"}
	fmt.Println("Closing target_add__run...")
	templates.ExecuteTemplate(w, "layout", infos)
}

func target_delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for target_delete...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/target_delete.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func target_delete__run(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting target_delete__run...")
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	target := data.Target{
		Url: r.PostFormValue("url"),
	}
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
	targets, err := user.UsersTargetsByUser()
	targetToBeDeleted, err := user.UsersTargetsByUserAndUrl(target.Url)
	if (data.Target{}) == targetToBeDeleted {
		// If the target inserted by the user does not exists
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/private.navigation.html",
				"templates/targets.html"))
		type TempStruct struct {
			User    data.User
			Targets []data.Target
			Message string
		}
		infos := TempStruct{user, targets, "Target does not exists"}
		templates.ExecuteTemplate(w, "layout", infos)
	} else {
		// If the target inserted by the user exists
		err := targetToBeDeleted.DeleteUserTargetByUserAndTarget(user)
		if err != nil {
			danger(err, "Cannot delete user <--> target")
		} else {
			templates := template.Must(
				template.ParseFiles(
					"templates/layout.html",
					"templates/private.navigation.html",
					"templates/targets.html"))
			type TempStruct struct {
				User    data.User
				Targets []data.Target
				Message string
			}
			infos := TempStruct{user, targets, "Target successfully deleted"}
			templates.ExecuteTemplate(w, "layout", infos)
		}
	}
}
