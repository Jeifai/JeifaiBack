package main

import (
	"encoding/json"
	"fmt"
	"html/template"
    "net/http"

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
    
    fmt.Println(user)
    
    _ = err

}
