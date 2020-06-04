package main

import (
	"fmt"
	"html/template"
	"net/http"

	// "io/ioutil"
	"encoding/json"

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
			"templates/private.navigation.html",
			"templates/keywords.html"))

	type TempStruct struct {
		User    data.User
		Targets []string
	}

	struct_targets, err := user.UsersTargetsByUser()

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Url)
	}

	infos := TempStruct{user, arr_targets}
	templates.ExecuteTemplate(w, "layout", infos)

	_ = err
}

func putKeywordsTargets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting putKeywordsTargets...")
	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		panic(err.Error())
	}

	type TempResponse struct {
		SelectedTargets []string     `json:"selectedTargets"`
		Keyword         data.Keyword `json:"newKeyword"`
	}

	response := TempResponse{}

	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		panic(err.Error())
	}

	// Before creating the relation user <-> target, check if it is not already present
	err = response.Keyword.KeywordByText()
	if err != nil {
		panic(err.Error())
	}

	// If keyword does not exist, create it
	if response.Keyword.Id == 0 {
		response.Keyword.CreateKeyword()
	}

	// Get target's detail based on
	targets, err := data.TargetsByUrls(response.SelectedTargets)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("response.Keyword")
	fmt.Println(response.Keyword)
	fmt.Println("targets")
	fmt.Println(targets)
	fmt.Println("user")
	fmt.Println(user)

	data.SetUserTargetKeyword(user, targets, response.Keyword)

	_ = err

	// if the relation userid, targetid, keywordid does not exist, create it, otherwise return a proper message
}
