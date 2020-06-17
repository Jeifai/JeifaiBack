package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"

	"github.com/go-playground/validator"

	"./data"
)

func targets(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
            "templates/topbar.html",
			"templates/sidebar.html",
			"templates/targets.html"))

	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		User data.User
		Data []data.Target
	}

	targets, err := user.UsersTargetsByUser()

	infos := TempStruct{user, targets}
	templates.ExecuteTemplate(w, "layout", infos)
}

func putTarget(w http.ResponseWriter, r *http.Request) {
	var target data.Target
	err := json.NewDecoder(r.Body).Decode(&target)

	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	url_parsed, err := url.Parse(target.Url)
	target.Host = url_parsed.Host

	validate := validator.New()
	err = validate.Struct(target)

	var message string
	added := false

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "Url" {
				if err.Tag() == "url" {
					temp_message = `The URL inserted is not valid`
				}
				message = red_1 + temp_message + red_2
			}
		}
	}

	if len(message) == 0 {

		// Try to create a target
		if err := target.CreateTarget(); err != nil {
			// If already exists, get its url
			err := target.TargetByUrl()
			if err != nil {
				panic(err.Error())
			}
		}

		// Before creating the relation user <-> target, check if it is not already present
		_, err := user.UsersTargetsByUserAndUrl(target.Url)

		if err != nil {

			// If the relation does not exists create a new relation
			target.CreateUserTarget(user)

			green_1 := `<p style="color:green">`
			green_2 := `</p>`
			message = green_1 + "Target successfully added" + green_2
			added = true
		} else {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			message = red_1 + "Target already exists" + red_2
		}
	}

	type TempStruct struct {
		Message string
		Added   bool
	}

	infos := TempStruct{message, added}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func removeTarget(w http.ResponseWriter, r *http.Request) {
	var target data.Target
	err := json.NewDecoder(r.Body).Decode(&target)

	sess, err := session(r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		panic(err.Error())
	}

	url_parsed, err := url.Parse(target.Url)
	target.Host = url_parsed.Host

	// Get the target to delete
	target, err = user.UsersTargetsByUserAndUrl(target.Url)
	if err != nil {
		panic(err.Error())
	}

	// Fill Deleted_At
	err = target.SetDeletedAtInUsersTargetsByUserAndTarget(user)
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct{ Message string }
	infos := TempStruct{"Target successfully removed"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
