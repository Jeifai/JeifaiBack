package main

import (
	"fmt"
	"html/template"
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

	r.HandleFunc("/targets", targets).Methods("GET")
	r.HandleFunc("/targets", putTarget).Methods("PUT")
	r.HandleFunc("/targets/remove", removeTarget).Methods("PUT")
	r.HandleFunc("/targets/all", targetsAll).Methods("GET")

	r.HandleFunc("/results", results)

	r.HandleFunc("/test", test)
	r.HandleFunc("/test_side", test_side)

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

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting test...")
	test_template := template.Must(template.ParseFiles("templates/test.html"))
	fmt.Println("Closing test...")
	test_template.ExecuteTemplate(w, "test.html", nil)
}

func test_side(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting test_side...")
	test_side_template := template.Must(template.ParseFiles("templates/test_side.html"))
	fmt.Println("Closing test_side...")
	test_side_template.ExecuteTemplate(w, "test_side.html", nil)
}
