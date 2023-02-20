package ukraine_scraper

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/enkemmc/go_tools"
	"github.com/enkemmc/notification_app"
)

const APP_NAME_HEADER = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36 Edg/99.0.1150.46"

var STRING_HITS []string

type UkraineScraper struct {
	exitChannel chan bool
	urlsChannel chan []*notification_app.UrlData
	name        string
}

func (this UkraineScraper) GetExitChannel() chan bool {
	return this.exitChannel
}
func (this UkraineScraper) GetUrlsChannel() chan []*notification_app.UrlData {
	return this.urlsChannel
}
func (this UkraineScraper) GetName() string {
	return this.name
}

func StartUkraineScraper() notification_app.LinkProvider {
	if bs, err := os.ReadFile("ukraine_scraper/terms.json"); err == nil {
		err = json.Unmarshal(bs, &STRING_HITS)
	} else {
		log.Fatal(err)
	}
	urlsChannel, exitChannel := startFetchLoop()
	name := "UATV Youtube"
	return UkraineScraper{
		urlsChannel,
		exitChannel,
		name,
	}
}

func startFetchLoop() (chan bool, chan []*notification_app.UrlData) {
	defaultDuration := 60 * time.Second
	ticker := time.NewTicker(defaultDuration)
	done := make(chan bool)

	urlsChan := make(chan []*notification_app.UrlData, 10) // interestingly, if we dont set a length, this will block
	tickImmediately := make(chan bool)
	fetchAndSend(urlsChan, tickImmediately)

	go func(urlsChan chan []*notification_app.UrlData, tickNow chan bool) {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fetchAndSend(urlsChan, tickNow)
			case <-tickNow:
				go fetchAndSend(urlsChan, tickNow)
				ticker.Reset(defaultDuration)
			}
		}
	}(urlsChan, tickImmediately)

	return done, urlsChan
}

func fetchAndSend(urlsChan chan []*notification_app.UrlData, tickNow chan bool) {
	urlsMap := fetchAndRead(tickNow)
	urls := []*notification_app.UrlData{}
	for _, entry := range urlsMap {
		urls = append(urls, entry)
	}
	urlsChan <- urls
}

func fetchAndRead(tickNow chan bool) map[string]*notification_app.UrlData {
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=UCOmfcmDrWs7iJrXx7V5Cnwg"
	client := http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = http.Header{
		"Host":       []string{"www.youtube.com"},
		"User-agent": []string{APP_NAME_HEADER},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	set := make(map[string]*notification_app.UrlData) // this set will contain the urls to any yt videos
	// check server status here
	if res.StatusCode == 429 {
		go_tools.PrintWithTimestamp("returned a 429 code\nretrying in 5 seconds")
		go func(tickNow chan bool) {
			time.Sleep(5 * time.Second)
			tickNow <- true
		}(tickNow)
		return set
	} else {
		go_tools.PrintWithTimestamp("returned a 200 code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(feed.Entries); i++ {
		var lp notification_app.UrlData = feed.Entries[i]
		set[lp.GetTitle()] = &lp
	}

	return set
}
