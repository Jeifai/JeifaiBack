package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	DbConnect()
	defer Db.Close()

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}

    notifications, err := GetNotifications(scrapers)
	if err != nil {
		panic(err.Error())
    }
    
    emails := GetEmails(notifications)
	fmt.Println(emails)

	file, _ := json.MarshalIndent(emails, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0o644)
}
