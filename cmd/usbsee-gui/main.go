package main

//go:generate go-localize -input ../../internal/localizations/src -output ../../internal/localizations

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/asaskevich/EventBus"
	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/oldgui"
)

//go:embed logo.png
var logo []byte

var topWindow *oldgui.TopWindow
var deviceChooser *oldgui.DeviceChooser

func makeLogo() fyne.Resource {
	return fyne.NewStaticResource("Usbsee Logo", logo)
}

func main() {

	l := localizations.New("en", "en")
	bus := EventBus.New()

	a := app.NewWithID("com.blacktau.usbsee")

	a.SetIcon(makeLogo())

	topWindow = oldgui.MakeTopWindow(a, l, &bus)
	deviceChooser = oldgui.MakeDeviceChooser(a, l, &bus)
	topWindow.ShowAndRun()
}
