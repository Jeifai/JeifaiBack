package main

import (
	"./data"
    "fmt"
	"html/template"
    "net/http"
)

func jobs(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Generating HTML for jobs...")
    sess, err := session(w, r)
    user, err := data.UserByEmail(sess.Email)
	if err != nil {
		danger(err, "Cannot find user")
	}
    jobs, err := user.JobsByUser()
	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/jobs.html"))
	type TempStruct struct {
		User    data.User
		Jobs    []data.Job
		Message string
	}
    infos := TempStruct{user, jobs, "Here the list of all your job opportunities"}
	templates.ExecuteTemplate(w, "layout", infos)
}