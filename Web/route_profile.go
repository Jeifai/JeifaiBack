package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"
)

func profile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for profile...")
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/topbar.html",
			"templates/sidebar.html",
			"templates/profile.html"))

	sess, err := session(r)
	user, err := UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	type PublicUser struct {
		UserName    string
		Email       string
		FirstName   sql.NullString
		LastName    sql.NullString
		DateOfBirth sql.NullString
		Country     sql.NullString
		City        sql.NullString
		Gender      sql.NullString
	}

	type TempStruct struct {
		User PublicUser
	}

	var publicUser PublicUser
	publicUser.UserName = user.UserName
	publicUser.Email = user.Email
	publicUser.FirstName.String = user.FirstName.String
	publicUser.LastName.String = user.LastName.String
	publicUser.Country.String = user.Country.String
	publicUser.City.String = user.City.String
	publicUser.DateOfBirth.String = user.DateOfBirth.String
	publicUser.Gender.String = user.Gender.String

	infos := TempStruct{publicUser}
	templates.ExecuteTemplate(w, "layout", infos)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting updateProfile...")

	sess, err := session(r)
	user, err := UserById(sess.UserId)

	err = json.NewDecoder(r.Body).Decode(&user)

	if user.CurrentPassword != "" {
		user.CurrentPassword = Encrypt(user.CurrentPassword)
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
			if err.Field() == "UserName" {
				temp_message = `UserName cannot be empty`
				messages = append(messages, red_1+temp_message+red_2)
			}
		}
	}

	if len(messages) == 0 {
		// Query will always update the password
		if user.NewPassword != "" { // User wants to change password
			user.NewPassword = Encrypt(user.NewPassword)
		} else { // User does not want to change the password
			user.NewPassword = user.CurrentPassword
		}

		fmt.Println("Updating user infos...")
		user.UpdateUser()
		user.UpdateUserUpdates()

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
