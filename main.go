package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {

	r := mux.NewRouter()
	files := http.FileServer(http.Dir(config.Static))
	r.Handle("/static/", http.StripPrefix("/static/", files))

	r.HandleFunc("/", index)

	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/signup_account", signupAccount)
	r.HandleFunc("/authenticate", authenticate)

	r.HandleFunc("/targets", targets).Methods("GET")
	r.HandleFunc("/targets", putTarget).Methods("PUT")
	r.HandleFunc("/targets/{url}", deleteTarget).Methods("DELETE")
	r.HandleFunc("/targets/all", targetsAll).Methods("GET")

	r.HandleFunc("/results", results)

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
