package main

import (
    "./data"
    "fmt"
	"net/http"
)

func targets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for targets...")
	generateHTML(w, nil, "layout", "private.navbar", "targets")
}

func target_add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for target_add...")
	generateHTML(w, nil, "layout", "private.navbar", "target_add")
}

func target_save(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting target_add_url...")
    err := r.ParseForm()
    sess, err := session(w, r)
	if err != nil {
		danger(err, "Cannot parse form")
	}
	target := data.Target{
		Url: r.PostFormValue("url"),
    }
	if err := target.CreateTarget(); err != nil {
		danger(err, "Cannot create target")
	}
	user, err := data.UserByEmail(sess.Email)
    if err != nil {
        danger(err, "Cannot find user")
    }
	if err := target.CreateUserTarget(user); err != nil {
		danger(err, "Cannot create relation user <--> target")
	}
	fmt.Println("Closing target_add_url...")
	http.Redirect(w, r, "/targets", 302)
}
