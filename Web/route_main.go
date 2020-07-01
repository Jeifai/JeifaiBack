package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting index...")
	sess, err := session(r)
	if err != nil {
		fmt.Println("Generating HTML for index, user not logged in...")
		/**
				templates := template.Must(
					template.ParseFiles(
						"templates/layout.html",
						"templates/public.navigation.html",
		                "templates/index.html"))*/
		templates := template.Must(template.ParseFiles("templates/tempIndex.html"))
		templates.ExecuteTemplate(w, "layout", nil)
	} else {
		user, err := UserById(sess.UserId)
		if err != nil {
			danger(err, "Cannot find user")
		}
		fmt.Println("Generating HTML for index, user logged in...")
		templates := template.Must(
			template.ParseFiles(
				"templates/layout.html",
				"templates/topbar.html",
				"templates/sidebar.html",
				"templates/index.html"))
		type TempStruct struct {
			User User
		}
		infos := TempStruct{user}
		templates.ExecuteTemplate(w, "layout", infos)
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	login_template := template.Must(template.ParseFiles("templates/test.html"))
	login_template.ExecuteTemplate(w, "test.html", nil)
}
