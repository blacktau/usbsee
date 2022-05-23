package gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/logging"
)

type Gui struct {
	mainWindow    *gtk.ApplicationWindow
	actionEntries []ActionEntry
	application   *gtk.Application
	i18n          *localizations.Localizer
	logger        logging.Logger
	deviceChooser DeviceChooser
}

func MakeGui(application *gtk.Application, l *localizations.Localizer, logger logging.Logger) (*Gui, error) {

	gui := &Gui{
		application: application,
		i18n:        l,
		logger:      logger,
	}

	mainWindow, err := makeAppWindow(application)

	if err != nil {
		return nil, fmt.Errorf("failed to create main window: %w", err)
	}

	gui.mainWindow = mainWindow

	deviceChooser, err := makeDeviceChooser(gui.mainWindow, logger)

	if err != nil {
		gui.logger.Errorf("Failed to create device chooser: %v", err)
		os.Exit(-1)
	}

	gui.deviceChooser = *deviceChooser
	application.AddWindow(deviceChooser.dialog)

	gui.mapActionEntries(application)

	logger.Debugf("deviceChooser: %v", deviceChooser)
	logger.Debugf("gui.deviceChooser: %v", gui.deviceChooser)

	return gui, nil
}

func (g Gui) Start() {
	g.logger.Debugf("Starting! %v", g.mainWindow)
	g.mainWindow.ShowAll()
}

func (g Gui) mapActionEntries(application *gtk.Application) {
	actionEntries := []ActionEntry{
		{
			Name:       "new",
			Enabled:    true,
			OnActivate: func() { g.newSession() },
		},
		{
			Name:       "open",
			Enabled:    true,
			OnActivate: func() { g.openSession() },
		},
		{
			Name:       "save",
			Enabled:    false,
			OnActivate: func() { g.saveSession() },
		},
		{
			Name:       "save-as",
			Enabled:    false,
			OnActivate: func() { g.saveSessionAs() },
		},
		{
			Name:       "device-selected",
			Enabled:    true,
			OnActivate: func() { g.deviceSelected() },
		},
	}

	mapActionEntries(application, actionEntries)
}

func (g Gui) newSession() {
	g.logger.Debug("NEW SESSION!!!")
	g.deviceChooser.Show()
}

func (g Gui) openSession() {
	g.logger.Debug("OPEN SESSION!!!")
}

func (g Gui) saveSession() {
	g.logger.Debug("SAVE SESSION!!!")
}

func (g Gui) saveSessionAs() {
	g.logger.Debug("SAVE SESSION AS!!!")
}

func (g Gui) deviceSelected() {
	g.logger.Debug("DEVICE SELECTED AS!!!")
}

func isApplicationWindow(obj glib.IObject) (*gtk.ApplicationWindow, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.ApplicationWindow); ok {
		return win, nil
	}

	return nil, fmt.Errorf("not a *gtk.ApplicationWindow")
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}

	return nil, fmt.Errorf("not a *gtk.Dialog")
}

func isListStore(obj glib.IObject) (*gtk.ListStore, error) {
	if store, ok := obj.(*gtk.ListStore); ok {
		return store, nil
	}

	return nil, fmt.Errorf("not a *gtk.ListStore")
}

func isButton(obj glib.IObject) (*gtk.Button, error) {
	if button, ok := obj.(*gtk.Button); ok {
		return button, nil
	}

	return nil, fmt.Errorf("not a *gtk.Button")
}
