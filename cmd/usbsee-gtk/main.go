package main

//go:generate go-localize -input ../../internal/localizations/src -output ../../internal/localizations

import (
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/blacktau/usbsee/internal/gui"
	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/logging"
)

const appID = "com.github.blacktau.usbsee"

func main() {

	logger := logging.MakeSugaredLogger()

	localizer := localizations.New("en", "en")

	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)

	if err != nil {
		logger.Fatalf("Failed to create application", err)
		os.Exit(1)
	}

	application.Connect("activate", func() {
		var err error

		ui, err := gui.MakeGui(application, localizer, logger)

		if err != nil {
			logger.Fatalf("Failed to start main window: %v", err)
			os.Exit(1)
		}

		ui.Start()
	})

	os.Exit(application.Run(os.Args))
}
