package main

import (
	"github.com/enkemmc/notification_app/notification_app"
	"github.com/enkemmc/notification_app/scraper"
	"github.com/enkemmc/notification_app/tools"
)

func main() {
	tools.PrintWithTimestamp("starting")
	provider := scraper.StartRedditScraper()
	app := notification_app.NewNotificationApp()
	app.AddTopic(provider)
	app.Start()
}
