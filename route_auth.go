package main

import (
	"./data"
	"fmt"
	"html/template"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting login...")
    login_template := template.Must(template.ParseFiles("templates/login.html"))
	fmt.Println("Closing login...")
	login_template.ExecuteTemplate(w, "login.html", nil)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting authenticate...")
	err := r.ParseForm()
	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
        http.SetCookie(w, &cookie)
	    fmt.Println("Closing login...")
		http.Redirect(w, r, "/", 302)
	} else {
        fmt.Println("Log in not valid...")
	    fmt.Println("Closing login...")
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting logout...")
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	fmt.Println("Closing logout...")
	http.Redirect(w, r, "/", 302)
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for signup...") 
	generateHTML(w, nil, "layout", "public.navbar", "signup")
}

func signupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting signupAccount...")
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
    }
	fmt.Println("Closing signupAccount...")
	http.Redirect(w, r, "/login", 302)
}
