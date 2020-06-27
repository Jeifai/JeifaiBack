package main

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
}
