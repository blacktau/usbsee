package main

//go:generate go-localize -input ../../internal/localizations/src -output ../../internal/localizations

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/blacktau/usbsee/internal/gui"
	"github.com/blacktau/usbsee/internal/localizations"
	"golang.org/x/text/message"
)

//go:embed logo.png
var logo []byte

var topWindow fyne.Window
var printer message.Printer

func makeToolBar() *widget.Toolbar {

	return widget.NewToolbar()
}

func makeUI() *widget.Toolbar {

	top := makeToolBar()

	return top
}

func makeLogo() fyne.Resource {
	return fyne.NewStaticResource("Usbsee Logo", logo)
}

func main() {

	l := localizations.New("en", "en")

	a := app.NewWithID("com.blacktau.usbsee")

	a.SetIcon(makeLogo())

	topWindow = a.NewWindow("Usbsee")
	topWindow.Resize(fyne.NewSize(100.0, 100.0))
	topWindow.SetMainMenu(gui.MakeMainMenu(&a, l))
	topWindow.SetMaster()

	gui.MakeDeviceChooser(a, l)

	topWindow.SetContent(container.NewBorder(makeUI(), nil, nil, nil))
	topWindow.ShowAndRun()
}
