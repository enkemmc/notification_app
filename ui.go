package notification_app

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	accordion *widget.Accordion
}

func StartUI(appid string) (*widget.Accordion, *fyne.Window, *fyne.App) {
	app := app.NewWithID(appid)
	window := app.NewWindow("some window name")
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

func BuildNewUrlWrapper(urlString string, vbox *fyne.Container, openURLfunc func(url *url.URL)) (fyne.CanvasObject, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	} else {
		hbox := container.NewHBox()
		hbox.Add(
			widget.NewLabel(urlString),
		)
		hbox.Add(
			widget.NewButton("Open", func() {
				openURLfunc(parsedUrl)
			}),
		)
		hbox.Add(
			widget.NewButton("remove", func() {
				vbox.Remove(hbox)
			}),
		)
		return hbox, nil
	}
}
