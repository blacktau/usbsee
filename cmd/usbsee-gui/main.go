package main

//go:generate go-localize -input ../../internal/localizations/src -output ../../internal/localizations

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/asaskevich/EventBus"
	"github.com/blacktau/usbsee/internal/gui"
	"github.com/blacktau/usbsee/internal/localizations"
)

//go:embed logo.png
var logo []byte

var topWindow *gui.TopWindow
var deviceChooser *gui.DeviceChooser

func makeLogo() fyne.Resource {
	return fyne.NewStaticResource("Usbsee Logo", logo)
}

func main() {

	l := localizations.New("en", "en")
	bus := EventBus.New()

	a := app.NewWithID("com.blacktau.usbsee")

	a.SetIcon(makeLogo())

	topWindow = gui.MakeTopWindow(a, l, &bus)
	deviceChooser = gui.MakeDeviceChooser(a, l, &bus)
	topWindow.ShowAndRun()
}
