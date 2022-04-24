package main

import (
	"github.com/enkemmc/notification_app"
	"github.com/enkemmc/use_not_app/reddit_scraper"
	ukrainescraper "github.com/enkemmc/use_not_app/ukraine_scraper"
)

func main() {
	redditProvider := reddit_scraper.StartRedditScraper()
	ukraineProvider := ukrainescraper.StartUkraineScraper()
	app := notification_app.NewNotificationApp("app_name")
	app.AddTopic(redditProvider)
	app.AddTopic(ukraineProvider)
	app.Start()

}
