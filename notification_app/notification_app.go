package notification_app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/enkemmc/notification_app/ui"
)

func NewNotificationApp() *NotificationApp {
	// something with ui?
	accordion, window := ui.StartUI()
	return &NotificationApp{
		data:      make(map[string]*TopicData),
		accordion: accordion,
		window:    window,
	}
}

type NotificationApp struct {
	data      map[string]*TopicData
	accordion *widget.Accordion
	window    *fyne.Window
}

func (app *NotificationApp) AddTopic(provider LinkProvider) {
	ai := ui.BuildNewAccordionItem(provider.GetName())
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
				app.RefreshUrls(urls, provider.GetName())
			}
		}
	}(ai.Detail.(*fyne.Container))
	//somewhere in here you need to start a goroutine that sends urlchannel data to this topic's accordionitem
	//for {
	//	select {
	//	case urls := <-provider.GetUrlsChannel():
	//		for i, url := range urls {
	//			fmt.Printf("%d %s\n", i, url)
	//		}
	//	case <-provider.GetExitChannel():
	//		break
	//	}
	//}
	app.data[provider.GetName()] = &td
}

func (app *NotificationApp) RefreshUrls(urls []string, topic string) {
	td := app.data[topic]
	vbox := (*td).accordionItem.Detail.(*fyne.Container)
	//	compare to whats in the urls list first, diff it against what should be visible in the vbox
	for _, url := range urls {
		if _, ok := td.urls[url]; !ok {
			// this url is new
			td.urls[url] = true
			row := ui.BuildNewUrlWrapper(url, vbox)
			vbox.Add(row)
			// need to store a reference to this row somewhere so we can call remove when "hide" is clicked
		}
	}
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

func start() {
}

type LinkProvider interface {
	GetExitChannel() chan bool
	GetUrlsChannel() chan []string
	GetName() string
}
