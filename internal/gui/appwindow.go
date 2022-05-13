package gui

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/blacktau/usbsee/internal/localizations"
)

//go:embed assets/logo.png
var logo []byte

var logoBuf *gdk.Pixbuf

func MakeAppWindow(application *gtk.Application, l *localizations.Localizer) (*gtk.ApplicationWindow, error) {

	appWindow, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		return nil, fmt.Errorf("Could not spawn root application window: %w", err)
	}

	appWindow.SetTitle(l.Get("usbsee.title"))

	appWindow.SetDefaultSize(1024, 768)

	logoBuf, err := loadAppIcon()
	if err != nil {
		log.Printf("Failed to load app icon: %v\n", err)
	}

	appWindow.SetIcon(logoBuf)

	body, err := buildMainBody(l)
	if err != nil {
		return nil, fmt.Errorf("Failed to create application UI: %w", err)
	}

	appWindow.Add(body)

	header, err := buildHeaderBar(l)

	if err != nil {
		return nil, fmt.Errorf("Failed to build header bar: %w", err)
	}

	appWindow.SetTitlebar(header)

	newAction := glib.SimpleActionNew("new", nil)
	newAction.Connect("activate", func() {
		fmt.Println("NEW!!!")
	})
	appWindow.AddAction(newAction)

	return appWindow, nil
}

func loadAppIcon() (*gdk.Pixbuf, error) {
	return gdk.PixbufNewFromBytesOnly(logo)
}
