package notification_app

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewNotificationApp(appid string) *NotificationApp {
	accordion, window, app := StartUI(appid)
	return &NotificationApp{
		data:      make(map[string]*TopicData),
		accordion: accordion,
		window:    window,
		app:       app,
	}
}

type NotificationApp struct {
	data      map[string]*TopicData
	accordion *widget.Accordion
	window    *fyne.Window
	app       *fyne.App
}

func (app *NotificationApp) AddTopic(provider LinkProvider) {
	ai := BuildNewAccordionItem(provider.GetName())
	td := TopicData{
		urls:           make(map[string]bool),
		accordionIndex: app.getNewIndex(),
		accordionItem:  ai,
	}
	app.accordion.Append(ai)
	go func(vbox *fyne.Container) {
		for {
			select {
			case urls := <-provider.GetUrlsChannel():
				app.refreshUrls(urls, provider.GetName())
			case <-provider.GetExitChannel():
				break
			}
		}
	}(ai.Detail.(*fyne.Container))
	app.data[provider.GetName()] = &td
}

func (app *NotificationApp) refreshUrls(urlDatas []*UrlData, topic string) {
	td := app.data[topic]
	vbox := (*td).accordionItem.Detail.(*fyne.Container)
	changeCount := 0
	for _, urldata := range urlDatas {
		url := (*urldata).GetUrl()
		if _, ok := td.urls[url]; !ok {
			td.urls[url] = true
			row, err := BuildNewUrlWrapper(urldata, vbox, app.openURL)
			if err == nil {
				vbox.Add(row)
				changeCount++
			}
		}
	}
	if changeCount > 0 {
		app.notify(changeCount)
	}
}
func (app *NotificationApp) notify(changes int) {
	(*app.app).SendNotification(fyne.NewNotification(fmt.Sprintf("%d new updates", changes), ""))
}

func (app *NotificationApp) openURL(urlString *url.URL) {
	(*app.app).OpenURL(urlString)
}

func (app *NotificationApp) Start() {
	(*app.window).ShowAndRun()
}

func (app *NotificationApp) getNewIndex() int {
	return len(app.data)
}

type TopicData struct {
	urls           map[string]bool
	accordionIndex int
	accordionItem  *widget.AccordionItem
}

type LinkProvider interface {
	GetExitChannel() chan bool
	GetUrlsChannel() chan []*UrlData
	GetName() string
}

// use this to store url and timestamp rather than just the url string
type UrlData interface {
	GetUrl() string
	GetElapsedTime() string
}
