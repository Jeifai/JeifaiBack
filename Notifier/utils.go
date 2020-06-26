package main

import "fmt"

func GetEmails(notifications []Notification) (emails []Email) {
    fmt.Println("Starting GetEmails...")
	var users []string
	for _, notif_1 := range notifications {
		if !Contains(users, notif_1.UserName) {
			var email Email
			email.UserName = notif_1.UserName
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

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}