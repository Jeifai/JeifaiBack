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
	user, err := data.UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	type PublicUser struct {
		UserName    string
		Email       string
		FirstName   string
		LastName    string
		DateOfBirth string
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
	publicUser.DateOfBirth = user.DateOfBirth
	publicUser.Gender = user.Gender

	infos := TempStruct{publicUser}
	templates.ExecuteTemplate(w, "layout", infos)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting updateProfile...")

	sess, err := session(r)
	user, err := data.UserById(sess.UserId)

	err = json.NewDecoder(r.Body).Decode(&user)

	if user.CurrentPassword != "" {
		user.CurrentPassword = data.Encrypt(user.CurrentPassword)
	}

	validate := validator.New()
	err = validate.Struct(user)

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
				temp_message = `New passwords do not match`
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
        fmt.Println("Updating user infos...")
		user.UpdateUser()
	}

	type TempStruct struct {
		Messages []string
	}
	infos := TempStruct{messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
