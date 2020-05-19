package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", index)

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	mux.HandleFunc("/targets", targets)
	mux.HandleFunc("/target_add", target_add)
	mux.HandleFunc("/target_add__run", target_add__run)
	mux.HandleFunc("/target_delete", target_delete)
	mux.HandleFunc("/target_delete__run", target_delete__run)

	mux.HandleFunc("/results", results)

	fmt.Println("Application is running")

	server := &http.Server{
		Addr:           "0.0.0.0:9090",
		Handler:        mux,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
    server.ListenAndServe()
}
