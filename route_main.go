package main

import (
	// "fmt"
	"./data"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public.navbar", "index")
	} else {
		user, err := data.UserByEmail(sess.Email)
		if err != nil {
			danger(err, "Cannot find user")
		}
		generateHTML(w, user, "layout", "private.navbar", "index")
	}
}
