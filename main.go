package main

import (
	"github.com/enkemmc/notification_app/notification_app"
	"github.com/enkemmc/notification_app/scraper"
)

func main() {
	provider := scraper.StartRedditScraper()
	app := notification_app.NewNotificationApp("app_name")
	app.AddTopic(provider)
	app.Start()
}
