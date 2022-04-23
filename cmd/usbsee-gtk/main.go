package main

//go:generate go-localize -input ../../internal/localizations/src -output ../../internal/localizations

import (
	"log"
	"os"

	"github.com/blacktau/usbsee/internal/gui"
	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appID = "com.github.blacktau.usbsee"

func main() {

	l := localizations.New("en", "en")

	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)

	if err != nil {
		log.Fatal("Failed to create application", err)
		os.Exit(1)
	}

	application.Connect("activate", func() { onActivate(application, l) })
	os.Exit(application.Run(os.Args))
}

func onActivate(application *gtk.Application, l *localizations.Localizer) {
	appWindow, err := gui.MakeAppWindow(application, l)

	if err != nil {
		log.Fatalf("Failed to start main window: %v", err)
		os.Exit(1)
	}

	appWindow.ShowAll()
}
