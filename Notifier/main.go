package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	DbConnect()
	defer Db.Close()

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}

	var notifications []Notification

	for _, elem := range scrapers {

		notifier := Notifier{1, 1, time.Now()}

		t_notifications, err := GetNotifications(notifier, elem.Id)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range t_notifications {
			notifications = append(notifications, elem)
		}
	}

	var emails []Email
	var users []string
	for _, notif_1 := range notifications {
		if !Contains(users, notif_1.UserName) {
			var email Email
			email.UserName = notif_1.UserName
			users = append(users, notif_1.UserName)
			var companies []string
			for _, notif_2 := range notifications {
				if notif_1.UserName == notif_2.UserName {
					if !Contains(companies, notif_2.CompanyName) {
						var company Company
						company.CompanyName = notif_2.CompanyName
						companies = append(companies, notif_2.CompanyName)
						var jobs []string
						for _, notif_3 := range notifications {
							if notif_1.UserName == notif_2.UserName {
								if notif_2.CompanyName == notif_3.CompanyName {
									if !Contains(jobs, notif_3.JobTitle) {
										var job Job
										job.JobTitle = notif_3.JobTitle
										job.JobUrl = notif_3.JobUrl
										company.Job = append(company.Job, job)
										jobs = append(jobs, notif_3.JobTitle)
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
	fmt.Println(emails)

	file, _ := json.MarshalIndent(emails, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0o644)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
