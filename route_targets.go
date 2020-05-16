package main

import (
	"./data"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type TargetOutput struct {
	User    data.User
	Targets []data.Target
	Message string
}

func targets(w http.ResponseWriter, r *http.Request) {
	info("Generating HTML for targets...")
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
	infos := TargetOutput{user, targets, ""}

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
	info("Starting target_add__run...")

	// Get useful information about user and session
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}

	// Get data provided by the user in the HTML form
	err = r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	temp_url := r.PostFormValue("url")
	temp_host, err := url.Parse(temp_url)

	// Instantiate a struct Target with all the data available atm
	target := data.Target{Url: temp_url, Host: temp_host.Host}

	// Try to create a target
	if err := target.CreateTarget(); err != nil {
		warning(err, "Cannot parse form")
		// If already exists, get its url
		err := target.TargetsByUrl()
		if err != nil {
			danger("Cannot retrive already existing Target")
		}
	}

	// Before creating the relation user <-> target, check if it is not already present
	targetsAlreadyExisting, err := user.UsersTargetsByUserAndUrl(target.Url)
	if (data.Target{}) == targetsAlreadyExisting {
		// If the relation does not exists create a new relation
		if err := target.CreateUserTarget(user); err != nil {
			danger(err, "Cannot create a new UsersTargets even if it doesn't exist")
		}
		// Finally extract all the relations of the user to properly print HTML
		targets, err := user.UsersTargetsByUser()
		if err != nil {
			danger(err, "Cannot retrive UsersTargets even if it was just created")
		}
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/private.navigation.html",
				"templates/targets.html"))
		infos := TargetOutput{user, targets, "Target successfully added"}
		info("Closing target_add__run...")
		templates.ExecuteTemplate(w, "layout", infos)
	} else {
		// If the relation already exists DO NOT create a new relation
		// Finally extract all the relations of the user to properly print HTML
		targets, err := user.UsersTargetsByUser()
		if err != nil {
			danger(err, "Cannot retrive UsersTargets")
		}
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/private.navigation.html",
				"templates/targets.html"))
		infos := TargetOutput{user, targets, "Target already present"}
		info("Closing target_add__run...")
		templates.ExecuteTemplate(w, "layout", infos)
	}
}

func target_delete(w http.ResponseWriter, r *http.Request) {
	info("Generating HTML for target_delete...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/target_delete.html"))
	templates.ExecuteTemplate(w, "layout", nil)
}

func target_delete__run(w http.ResponseWriter, r *http.Request) {
	info("Starting target_delete__run...")

	// Get useful information about user and session
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}

	// Get data provided by the user in the HTML form
	err = r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	temp_url := r.PostFormValue("url")
	temp_host, err := url.Parse(temp_url)

	// Instantiate a struct Target with all the data available atm
	target := data.Target{Url: temp_url, Host: temp_host.Host}

	// Based on the url provided by the user, understand if there is a target to delete
	targetToBeDeleted, err := user.UsersTargetsByUserAndUrl(target.Url)

	if (data.Target{}) == targetToBeDeleted {
		// If the target inserted by the user does not exists
		targets, err := user.UsersTargetsByUser()
		if err != nil {
			danger(err, "Cannot extract UsersTargets even if they should exist")
		}
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/private.navigation.html",
				"templates/targets.html"))
		infos := TargetOutput{user, targets, "Target does not exists"}
		templates.ExecuteTemplate(w, "layout", infos)
	} else {
		// If the target inserted by the user exists
		err := targetToBeDeleted.DeleteUserTargetByUserAndTarget(user)
		if err != nil {
			danger(err, "Cannot delete UserTarget")
		} else {
			targets, err := user.UsersTargetsByUser()
			if err != nil {
				danger(err, "Cannot extract UsersTargets even if they should exist")
			}
			templates := template.Must(
				template.ParseFiles(
					"templates/layout.html",
					"templates/private.navigation.html",
					"templates/targets.html"))
			infos := TargetOutput{user, targets, "Target successfully deleted"}
			templates.ExecuteTemplate(w, "layout", infos)
		}
	}
}
