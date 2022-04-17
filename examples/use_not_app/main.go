package main

import (
	"github.com/enkemmc/notification_app"
	"github.com/enkemmc/use_not_app/scraper"
)

func main() {
	provider := scraper.StartRedditScraper()
	app := notification_app.NewNotificationApp("app_name")
	app.AddTopic(provider)
	app.Start()
}
