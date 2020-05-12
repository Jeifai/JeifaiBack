package main

import (
	"fmt"
	"./data"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting index...")
	sess, err := session(w, r)
	if err != nil {
	    fmt.Println("Generating HTML for index, user not logged in...")   
		generateHTML(w, nil, "layout", "public.navbar", "index")
	} else {
		user, err := data.UserByEmail(sess.Email)
		if err != nil {
			danger(err, "Cannot find user")
        }
	    fmt.Println("Generating HTML for index, user logged in...") 
		generateHTML(w, user, "layout", "private.navbar", "index")
	}
}
