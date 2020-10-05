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

type Invitation struct {
	Email          string
	InvitationCode string
}

func main() {
	invitation := Invitation{
		"sadotix998@pickybuys.com",
		"d82a41ce-f5c7-4d0e-6845-7c56ffda00ad",
	}

	invitation.SendInvitationCode()
}

func (invitation *Invitation) SendInvitationCode() {
	fmt.Println(Gray(8-1, "Starting SendInvitationCode..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("templateInvitationCode.html")

	t, err = t.ParseFiles("templateInvitationCode.html")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(Blue("Sending email to -->"), Bold(Blue(invitation.Email)))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, invitation); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "robimalco@gmail.com")
	m.SetHeader("To", invitation.Email)
	m.SetHeader("Subject", "Hey here your invitation code!")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err.Error())
	}

	// SaveEmailIntoDb(invitation.Email, "SendInvitationCode")
}
