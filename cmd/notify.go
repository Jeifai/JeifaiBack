package cmd

import (
	"github.com/spf13/cobra"
)

var notifyUser string

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Run the notifier",
	Long:  `Run the notifier for specific users or for all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
		Notify(notifyUser)
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.Flags().StringVarP(&notifyUser, "notify", "n", "", "Specify a user or all of them")
}

func Notify(user string) {
	DbConnect()
	defer Db.Close()

	scrapers := GetScrapers()

	var notifications []Notification
	if user == "all" {
		notifications = GetNotifications(scrapers)
	} else {
		notifications = GetUserNotifications(scrapers, user)
	}

	RunNotifer(notifications)
}

func RunNotifer(notifications []Notification) {
	emails := CreateEmails(notifications)
	for _, email := range emails {
		var notifier Notifier
		notifier.StartNotifierSession(email.UserId)
		SaveUserNotifications(notifications, email, notifier)
		SendMatches(email)
		SaveEmail(email.UserEmail, "SendEmailMatches")
	}
}
