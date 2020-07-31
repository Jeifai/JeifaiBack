package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
	"gopkg.in/gomail.v2"
)

type TargetsReady struct {
	Email          string
	Targets        []string
}

func main() {
	targetsready := TargetsReady{
		"robimalco@gmail.com",
		[]string{"Paintgun", "test1", "test2"},
	}

	targetsready.SendTargetsReady()
}

func (targetsready *TargetsReady) SendTargetsReady() {
	fmt.Println(Gray(8-1, "Starting SendTargetsReady..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("templateTargetsReady.html")

	t, err = t.ParseFiles("templateTargetsReady.html")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(targetsready.Email)))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, targetsready); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", targetsready.Email)
	m.SetHeader("Subject", "Hey, there is something for you!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}
}
