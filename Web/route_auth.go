package main

import (
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
	user, err := UserByEmail(r.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}

	if user.Password == Encrypt(r.PostFormValue("password")) {
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
		session := Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	fmt.Println("Closing logout...")
	http.Redirect(w, r, "/", 302)
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting signup...")
	login_template := template.Must(template.ParseFiles("templates/signup.html"))
	fmt.Println("Closing signup...")
	login_template.ExecuteTemplate(w, "signup.html", nil)
}

func signupAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting signupAccount...")
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	user := User{
		UserName: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	fmt.Println("Closing signupAccount...")
	http.Redirect(w, r, "/login", 302)
}
