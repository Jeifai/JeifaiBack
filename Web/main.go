package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	files := http.FileServer(http.Dir(config.Static))
	s := http.StripPrefix("/static/", files)
	r.PathPrefix("/static/").Handler(s)

	r.HandleFunc("/", index)

	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/signup_account", signupAccount)
	r.HandleFunc("/authenticate", authenticate)

	r.HandleFunc("/profile", profile).Methods("GET")
	r.HandleFunc("/profile", updateProfile).Methods("PUT")

	r.HandleFunc("/targets", targets).Methods("GET")
	r.HandleFunc("/targets", putTarget).Methods("PUT")
	r.HandleFunc("/targets/remove", removeTarget).Methods("PUT")

	r.HandleFunc("/keywords", keywords).Methods("GET")
	r.HandleFunc("/keywords", putKeyword).Methods("PUT")
	r.HandleFunc("/keywords/remove", removeKeyword).Methods("PUT")

	r.HandleFunc("/results", results)

	r.HandleFunc("/test", test)

	fmt.Println("Application is running")

	server := &http.Server{
		Addr:           "0.0.0.0:9090",
		Handler:        r,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
