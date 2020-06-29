package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-playground/validator"

	"./data"
)

func keywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for keywords...")
	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/topbar.html",
			"templates/sidebar.html",
			"templates/keywords.html"))

	struct_targets, err := user.UsersTargetsByUser()

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Url)
	}

	utks, err := user.GetUserTargetKeyword()

	type TempStruct struct {
		User    data.User
		Targets []string
		Utks    []data.UserTargetKeyword
	}

	infos := TempStruct{user, arr_targets, utks}
	templates.ExecuteTemplate(w, "layout", infos)

	_ = err
}

func putKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting putKeyword...")
	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	type TempResponse struct {
		SelectedTargets []string     `json:"selectedTargets" validate:"required"`
		Keyword         data.Keyword `json:"newKeyword"`
	}

	response := TempResponse{}

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err.Error())
	}

	validate := validator.New()
	err = validate.Struct(response)

	var messages []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			red_1 := `<p style="color:red">`
			red_2 := `</p>`
			var temp_message string
			if err.Field() == "SelectedTargets" {
				if err.Tag() == "required" {
					temp_message = `Targets cannot be empty`
				}
			} else if err.Field() == "Text" {
				if err.Tag() == "required" {
					temp_message = `Keyword cannot be empty`
				}
				if err.Tag() == "min" {
					temp_message = `Keyword inserted is too short`
				}
				if err.Tag() == "max" {
					temp_message = `Keyword inserted is too long`
				}
			}
			messages = append(messages, red_1+temp_message+red_2)
		}
	}

	if len(messages) == 0 {

		// Before creating the relation user <-> target,
		// check if it is not already present
		err = response.Keyword.KeywordByText()

		// If keyword does not exist, create it
		if response.Keyword.Id == 0 {
			response.Keyword.CreateKeyword()
		}

		targets, err := data.TargetsByUrls(response.SelectedTargets)
		if err != nil {
			panic(err.Error())
		}

		err = data.SetUserTargetKeyword(user, targets, response.Keyword)
		if err != nil {
			panic(err.Error())
		}
		temp_message := `<p style="color:green">Successfully added</p>`
		messages = append(messages, temp_message)
	}

	var utks []data.UserTargetKeyword
	utks, err = user.GetUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct {
		Messages []string
		Utks     []data.UserTargetKeyword
	}

	infos := TempStruct{messages, utks}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}

func removeKeyword(w http.ResponseWriter, r *http.Request) {
	var utk data.UserTargetKeyword
	err := json.NewDecoder(r.Body).Decode(&utk)

	sess, err := session(r)
	user, err := data.UserByEmail(sess.Email)
	if err != nil {
		panic(err.Error())
	}

	target := data.Target{}
	target.Url = utk.TargetUrl
	err = target.TargetByUrl()
	if err != nil {
		panic(err.Error())
	}

	keyword := data.Keyword{}
	keyword.Text = utk.KeywordText
	err = keyword.KeywordByText()
	if err != nil {
		panic(err.Error())
	}

	utk.UserId = user.Id
	utk.TargetId = target.Id
	utk.KeywordId = keyword.Id

	err = utk.SetDeletedAtIntUserTargetKeyword()
	if err != nil {
		panic(err.Error())
	}

	type TempStruct struct{ Message string }
	infos := TempStruct{"Successfully removed"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(infos)
}
