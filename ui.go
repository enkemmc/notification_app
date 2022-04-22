package notification_app

import (
	"log"
	"net/url"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	accordion *widget.Accordion
}

const iconPath = "/resources/Susge.png"

func StartUI(appid string) (*widget.Accordion, *fyne.Window, *fyne.App) {
	app := app.NewWithID(appid)
	_, path, _, _ := runtime.Caller(0)
	path = filepath.Join(filepath.Dir(path), iconPath)
	r, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		log.Fatal(err)
	} else {
		app.SetIcon(r)
	}
	window := app.NewWindow("Notifications")
	content, accordion := buildContent()
	window.SetContent(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(content, nil, nil, nil),
		content,
	))
	window.Resize(fyne.Size{Width: 400, Height: 320})
	window.CenterOnScreen()
	return accordion, &window, &app
}

func buildContent() (*fyne.Container, *widget.Accordion) {
	accordion := widget.NewAccordion(
	// the topics are going to go here
	)
	content := fyne.NewContainerWithLayout(
		layout.NewCenterLayout(),
		accordion,
	)
	return content, accordion
}

func BuildNewAccordionItem(title string) *widget.AccordionItem {
	ai := widget.NewAccordionItem(
		title,
		container.NewVBox(),
	)
	return ai
}

func BuildNewUrlWrapper(urlData *UrlData, vbox *fyne.Container, openURLfunc func(url *url.URL)) (fyne.CanvasObject, error) {
	title := (*urlData).GetTitle()
	urlString := (*urlData).GetUrl()
	elapsedTime := (*urlData).GetElapsedTime()
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	} else {
		hbox := container.NewHBox()
		hbox.Add(
			widget.NewLabel(title),
		)
		hbox.Add(
			widget.NewLabel(elapsedTime),
		)
		hbox.Add(
			widget.NewButton("Open", func() {
				openURLfunc(parsedUrl)
			}),
		)
		hbox.Add(
			widget.NewButton("Clear", func() {
				vbox.Remove(hbox)
			}),
		)
		return hbox, nil
	}
}
