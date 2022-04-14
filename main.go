package main

import (
	"fmt"
	"time"

	"github.com/enkemmc/notification_app/scraper"
)

func main() {
	urlsChan := scraper.StartFetchLoop()
	exitChan := make(chan bool)
	go func(exitChan chan bool) {
		time.Sleep(60)
		exitChan <- true

	}(exitChan)
	for {
		select {
		case urls := <-urlsChan:
			for i, url := range urls {
				fmt.Printf("%d %s\n", i, url)
			}
		case <-exitChan:
			return
		}
	}
}
