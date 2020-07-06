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

func CreateEmailsStruct(notifications []Notification) (emails []Email) {
	fmt.Println(Gray(8-1, "Starting CreateEmailsStruct..."))
	var users []string
	for _, notif_1 := range notifications {
		if !Contains(users, notif_1.UserName) {
			var email Email
			email.UserId = notif_1.UserId
			email.UserName = notif_1.UserName
			email.UserEmail = notif_1.UserEmail
			users = append(users, notif_1.UserName)
			var companies []string
			for _, notif_2 := range notifications {
				if notif_1.UserName == notif_2.UserName {
					if !Contains(companies, notif_2.Name) {
						var company Company
						company.Name = notif_2.Name
						companies = append(companies, notif_2.Name)
						var jobs []string
						for _, notif_3 := range notifications {
							if notif_1.UserName == notif_2.UserName {
								if notif_2.Name == notif_3.Name {
									if !Contains(jobs, notif_3.Title) {
										var job Job
										job.MatchId = notif_3.MatchId
										job.Title = notif_3.Title
										job.Url = notif_3.Url
										company.Job = append(company.Job, job)
										jobs = append(jobs, notif_3.Title)
									}
								}
							}
						}
						email.Company = append(email.Company, company)
					}
				}
			}
			emails = append(emails, email)
		}
	}
	return
}

func SendEmails(emails []Email) {
	fmt.Println(Gray(8-1, "Starting SendEmails..."))
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	password := os.Getenv("PASSWORD")

	t := template.New("template.html")

	t, err = t.ParseFiles("template.html")
	if err != nil {
		panic(err.Error())
	}

	for _, email := range emails {

		fmt.Println(Blue("Sending email to -->"), Bold(Blue(email.UserEmail)))

		var notifier Notifier
		err := notifier.StartNotifierSession(email.UserId)
		if err != nil {
			panic(err.Error())
		}

		SaveNotification(notifier, email)

		var tpl bytes.Buffer
		if err := t.Execute(&tpl, email); err != nil {
			fmt.Println(err)
		}

		result := tpl.String()

		m := gomail.NewMessage()
		m.SetHeader("From", "robimalco@gmail.com")
		m.SetHeader("To", email.UserEmail)
		m.SetHeader("Subject", "Hello! There are new matches!")
		m.SetBody("text/html", result)

		d := gomail.NewDialer("smtp.gmail.com", 587, "robimalco@gmail.com", password)

		if err := d.DialAndSend(m); err != nil {
			panic(err.Error())
		}
	}
}
