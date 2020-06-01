package main

import (
	"encoding/json"
	"fmt"
	"html/template"
    "net/http"
    "time"

    "github.com/go-playground/validator"

	"./data"
)

type TempUser struct {
    UserName          string
    Email             string
    Password          string
    FirstName         string
    LastName          string
    DateOfBirth       time.Time
    Country           string
    City              string
    Gender            string
    CurrentPassword   string `validate:"required,eqfield=Password"`
    NewPassword       string `validate:"eqfield=RepeatNewPassword"`
    RepeatNewPassword string `validate:"eqfield=NewPassword"`
}
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

    var temp_user TempUser

    err = json.NewDecoder(r.Body).Decode(&temp_user)

    if temp_user.CurrentPassword != "" {
        temp_user.CurrentPassword = data.Encrypt(temp_user.CurrentPassword)
    }
    
    temp_user.Password = user.Password

    validate := validator.New()
	err = validate.Struct(temp_user)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
            fmt.Println(err.Field())
            var temp_message string
			if err.Field() == "CurrentPassword" {
                if err.Tag() == "required" {
                    temp_message = `<p style="color:red">Current password cannot be empty</p>`
                } else if err.Tag() == "eqfield" {
                    temp_message = `<p style="color:red">Wrong current password</p>`
                }
				messages = append(messages, temp_message)
			}
			if err.Field() == "NewPassword" {
				temp_message = `<p style="color:red">The new passwords do not match</p>`
				messages = append(messages, temp_message)
			}
		}
	}

	if len(messages) == 0 {
		temp_message := `<p style="color:green">Changes saved</p>`
		messages = append(messages, temp_message)
	}

	type TempStruct struct {
		Messages []string
	}
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
