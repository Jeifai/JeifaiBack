package main

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

	emails := CreateEmails(notifications)

	for _, email := range emails {

		var notifier Notifier
		err := notifier.StartNotifierSession(email.UserId)
		if err != nil {
			panic(err.Error())
		}

		SaveUserNotifications(notifications, email, notifier)
		if err != nil {
			panic(err.Error())
		}

		SendMatches(email)

		SaveEmail(email.UserEmail, "SendEmailMatches")
	}
}
