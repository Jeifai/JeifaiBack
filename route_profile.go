package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"

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

func updateProfile(w http.ResponseWriter, r *http.Request) {
	sess, err := session(r)
	user, err := data.UserByEmail(sess.Email)

	err = json.NewDecoder(r.Body).Decode(&user)

	user.CurrentPassword = data.Encrypt(user.CurrentPassword)

	validate := validator.New()

	err = validate.Struct(user)

	var errors []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "CurrentPassword" {
				temp_message := `<p style="color:red">Wrong current password</p>`
				errors = append(errors, temp_message)
			}
			if err.Field() == "NewPassword" {
				temp_message := `<p style="color:red">The new passwords do not match</p>`
				errors = append(errors, temp_message)
			}
		}
	}

	if len(errors) == 0 {
		temp_message := `<p style="color:green">Changes saved</p>`
		errors = append(errors, temp_message)
	}

	type TempStruct struct {
		Errors []string
	}
	infos := TempStruct{errors}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
