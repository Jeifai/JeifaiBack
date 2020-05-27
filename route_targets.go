package main

import (
	"./data"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
)

func targets(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/targets.html"))

	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	_ = err
	type TempStruct struct {
		User data.User
	}
	infos := TempStruct{user}
	templates.ExecuteTemplate(w, "layout", infos)
}

func targetsAll(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
	targets, err := user.UsersTargetsByUser()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(targets)
}

func putTarget(w http.ResponseWriter, r *http.Request) {
	var target data.Target
	err := json.NewDecoder(r.Body).Decode(&target)

	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}

	url_parsed, err := url.Parse(target.Url)
	target.Host = url_parsed.Host

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
		type TempStruct struct {
			Message string
			Added   bool
		}
		infos := TempStruct{"Target successfully added", true}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(infos)
	} else {
		type TempStruct struct {
			Message string
			Added   bool
		}
		infos := TempStruct{"Target already exists", false}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(infos)
	}
}

func deleteTarget(w http.ResponseWriter, r *http.Request) {

	info("Starting target_delete__run...")

	params := mux.Vars(r)

	// Get useful information about user and session
	sess, err := session(w, r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}

	// Instantiate a struct Target with all the data available atm
	target := data.Target{Url: params["url"]}

	// Based on the url provided by the user, understand if there is a target to delete
	targetToBeDeleted, err := user.UsersTargetsByUserAndUrl(target.Url)

	// If the target inserted by the user exists
	err = targetToBeDeleted.DeleteUserTargetByUserAndTarget(user)
	_ = err
}
