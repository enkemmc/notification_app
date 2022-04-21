package scraper

import (
	"encoding/xml"
	"log"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName   xml.Name `xml:"entry"`
	Id        string   `xml:"id"`
	Title     string   `xml:"title"`
	Content   string   `xml:"content"`
	Link      LinkTag  `xml:"link"`
	Published string   `xml:"published"`
}

func (e Entry) GetElapsedTime() string {
	start, err := time.Parse(time.RFC3339, e.Published)
	if err != nil {
		log.Fatal(err)
	}
	result := time.Since(start).Round(time.Second).String()
	return result
}

func (e Entry) GetTitle() string {
	return e.Title
}
func (e Entry) GetUrl() string {
	return e.Link.Href
}

type LinkTag struct {
	XMLNAME xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
}
