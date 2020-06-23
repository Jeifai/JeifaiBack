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

        // Prepare templates
        // Send emails
        // Save notifications into db
        users := []
		for _, elem := range notifications {
			fmt.Println(elem.CreatedAt)
			fmt.Println("\t", elem.UtkId)
		}
	}
}
