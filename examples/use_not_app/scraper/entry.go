package scraper

import "encoding/xml"

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
	Content string   `xml:"content"`
	Link    LinkTag  `xml:"link"`
}

type LinkTag struct {
	XMLNAME xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
}
