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

	type PublicUser struct {
		UserName    string
		Email       string
		FirstName   string
		LastName    string
		DateOfBirth time.Time
		Country     string
		City        string
		Gender      string
	}

	type TempStruct struct {
		User PublicUser
	}

	var publicUser PublicUser
	publicUser.UserName = user.UserName
	publicUser.Email = user.Email
	publicUser.FirstName = user.FirstName
	publicUser.LastName = user.LastName
	publicUser.Country = user.Country
	publicUser.City = user.City
	publicUser.Gender = user.Gender

	infos := TempStruct{publicUser}
	templates.ExecuteTemplate(w, "layout", infos)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	sess, err := session(r)
	user, err := data.UserByEmail(sess.Email)

	type TempUser struct {
		UserName          string
		Email             string `validate:"email"`
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
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "CurrentPassword" {
				if err.Tag() == "required" {
					temp_message = `Current password cannot be empty`
				} else if err.Tag() == "eqfield" {
					temp_message = `Wrong current password`
				}
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "NewPassword" {
				temp_message = `The new passwords do not match`
				messages = append(messages, red_1+temp_message+red_2)
			}
			if err.Field() == "Email" {
				temp_message = `Email inserted is not valid`
				messages = append(messages, red_1+temp_message+red_2)
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
