package gui

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	//go:embed assets/logo.png
	logo []byte

	//go:embed assets/main_window.glade
	mainWindowGlade string
)

func makeAppWindow(application *gtk.Application) (*gtk.ApplicationWindow, error) {

	builder, err := gtk.BuilderNewFromString(mainWindowGlade)

	if err != nil {
		return nil, fmt.Errorf("could not load ui defintion: %w", err)
	}

	winObj, err := builder.GetObject("main-window")
	if err != nil {
		return nil, fmt.Errorf("could find main-window in ui definition: %w", err)
	}

	appWindow, err := isApplicationWindow(winObj)
	if err != nil {
		return nil, fmt.Errorf("main-window is not an application window?!?: %w", err)
	}

	application.AddWindow(appWindow)

	logoBuf, err := loadAppIcon()
	if err != nil {
		log.Printf("Failed to load app icon: %v\n", err)
	}

	appWindow.SetIcon(logoBuf)

	return appWindow, nil
}

func loadAppIcon() (*gdk.Pixbuf, error) {
	return gdk.PixbufNewFromBytesOnly(logo)
}
