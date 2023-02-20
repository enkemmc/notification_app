package main

import (
	"github.com/enkemmc/notification_app"
	"github.com/enkemmc/use_not_app/reddit_scraper"
	"github.com/enkemmc/use_not_app/ukraine_scraper"
)

func main() {
	redditProvider := reddit_scraper.StartRedditScraper()
	ukraineProvider := ukraine_scraper.StartUkraineScraper()
	app := notification_app.NewNotificationApp("Notification App")
	app.AddTopic(redditProvider)
	app.AddTopic(ukraineProvider)
	app.Start()
}
