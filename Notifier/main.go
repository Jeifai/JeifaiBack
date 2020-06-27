package main

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
)

func main() {
	DbConnect()
	defer Db.Close()

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}

	notifications, err := PrepareNotifications(scrapers)
	if err != nil {
		panic(err.Error())
	}

	emails := CreateEmailsStruct(notifications)

	SendEmails(emails)

	file, _ := json.MarshalIndent(emails, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0o644)
}
