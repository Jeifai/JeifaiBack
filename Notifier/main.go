package main

import (
	"fmt"
	"time"
)

func main() {

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

        users := GetUniqueUsers(notifications)

        fmt.Println(users)

        /**

		for _, elem := range notifications {
            var email Email

            email.UserId = elem.UserId
            email.UserName = elem.UserName
            email.Company.CompanyName = elem.CompanyName
            email.Company.Job.JobTitle = elem.JobTitle
            email.Company.Job.JobUrl = elem.JobUrl

            emails = append(emails, email)
        }
        
        */
    }
    
    // fmt.Println(emails)
}

func GetUniqueUsers(notifications []Notification) (users []int) {
    for _, elem := range notifications {
        users = AppendIfMissing(users, elem.UserId)
    }
    return 
}

func AppendIfMissing(slice []int, element int) []int {
	for _, ele := range slice {
		if ele == element {
			return slice
		}
	}
	return append(slice, element)
}