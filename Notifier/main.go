package main

import (
	"fmt"
	"time"
)

func main() {

    var emails []Email

	DbConnect()
	defer Db.Close()

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {

		notifier := Notifier{1, 1, time.Now()}

		notifications, err := GetNotifications(notifier, elem.Id)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range notifications {
            var email Email
            var company Company
            var jobs []Job

            email.UserId = elem.UserId
            email.UserName = elem.UserName
            company.
            email.Companies

            emails = append(emails, email)
		}
    }
    
    fmt.Println(users)
}

func AppendIfMissing(slice []int, element int) []int {
	for _, ele := range slice {
		if ele == element {
			return slice
		}
	}
	return append(slice, element)
}