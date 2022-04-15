package scraper

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/enkemmc/notification_app/notification_app"
	"github.com/enkemmc/notification_app/tools"
)

const APP_NAME_HEADER = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36 Edg/99.0.1150.46"

// these are search terms we're looking for in posts
var STRING_HITS = [...]string{"code", "prerelease", "cutting"}

type RedditScraper struct {
	exitChan chan bool
	urlsChan chan []string
	name     string
}

func (rs *RedditScraper) GetUrlsChannel() chan []string {
	return rs.urlsChan
}
func (rs *RedditScraper) GetExitChannel() chan bool {
	return rs.exitChan
}
func (rs *RedditScraper) GetName() string {
	return rs.name
}

func StartRedditScraper() notification_app.LinkProvider {
	urlsChan, exitChan := startFetchLoop()
	name := "reddit_scraper"
	return &RedditScraper{
		exitChan,
		urlsChan,
		name,
	}
}

func startFetchLoop() (chan []string, chan bool) {
	tools.PrintWithTimestamp("starting fetch loop")
	defaultDuration := 30 * time.Second
	ticker := time.NewTicker(defaultDuration)
	done := make(chan bool)

	urlsChan := make(chan []string, 10) // interestingly, if we dont set a length, this will block
	tickImmediately := make(chan bool)
	fetchAndSend(urlsChan, tickImmediately)

	go func(urlsChan chan []string, tickNow chan bool) {
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

	return urlsChan, done
}

func fetchAndSend(urlsChan chan []string, tickNow chan bool) {
	urlsMap := fetchAndRead(tickNow)
	urls := []string{}
	for url, _ := range urlsMap {
		urls = append(urls, url)
	}
	urlsChan <- urls
}

// returns a set of imgPaths
func fetchAndRead(tickNow chan bool) map[string]bool {
	url := "https://www.reddit.com/r/MagicArena/new/.rss?sort=new"
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
		"Host":       []string{"www.reddit.com"},
		"User-agent": []string{APP_NAME_HEADER},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	set := make(map[string]bool) // this set will contain the urls to any images that match our conditions
	// check server status here
	if res.StatusCode == 429 {
		tools.PrintWithTimestamp("returned a 429 code\nretrying in 5 seconds")
		go func(tickNow chan bool) {
			time.Sleep(5 * time.Second)
			tickNow <- true
		}(tickNow)
		return set
	} else {
		tools.PrintWithTimestamp("returned a 200 code")
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
		searchEntry(feed.Entries[i], &set)
	}

	return set
}

func searchEntry(entry Entry, matchesSet *map[string]bool) {
	// re := regexp.MustCompile(`https://i.redd.it[^"]+`) // this is the regexp we're using to see if the post contains a link to an image in its body.  if it does, its a hit
	content := strings.ToLower(entry.Content)
	title := strings.ToLower(entry.Title)

	found := false

	// check to see if the content or title of this entry contain any of the terms we're searching for
	for _, s := range STRING_HITS {
		if found {
			break
		}

		if strings.Contains(content, s) || strings.Contains(title, s) {
			found = true
			// does the content have a link to an image?
			// res := re.Find([]byte(content))
			// if res != nil {
			// 	(*matchesSet)[string(res)] = true
			// }
			(*matchesSet)[entry.Link.Href] = true
		}
	}
}
