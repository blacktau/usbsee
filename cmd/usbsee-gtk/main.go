package main

//go:generate go-localize -input ../../internal/localizations/src -output ../../internal/localizations

import (
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/blacktau/usbsee/internal/gui"
	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/usb"
)

const appID = "com.github.blacktau.usbsee"

var actionEntries = []gui.ActionEntry{
	{
		Name:       "new",
		Enabled:    true,
		OnActivate: func() { newSession() },
	},
	{
		Name:       "open",
		Enabled:    true,
		OnActivate: func() { openSession() },
	},
	{
		Name:       "save",
		Enabled:    false,
		OnActivate: func() { saveSession() },
	},
	{
		Name:       "save-as",
		Enabled:    false,
		OnActivate: func() { saveSessionAs() },
	},
	// {
	// 	Name:       "device-selected",
	// 	Enabled:    true,
	// 	OnActivate: func() { deviceSelected() },
	// },
}

var localizer *localizations.Localizer
var mainWindow *gtk.ApplicationWindow
var logger *zap.SugaredLogger

func main() {

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	devLogger, _ := config.Build()
	defer devLogger.Sync()

	logger = devLogger.Sugar()

	localizer = localizations.New("en", "en")

	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)

	if err != nil {
		logger.Fatal("Failed to create application", err)
		os.Exit(1)
	}

	application.Connect("activate", func() { onActivate(application, localizer) })
	os.Exit(application.Run(os.Args))
}

func onActivate(application *gtk.Application, l *localizations.Localizer) {

	gui.MapActionEntries(application, actionEntries)

	var err error
	mainWindow, err = gui.MakeAppWindow(application, l)

	if err != nil {
		logger.Fatalf("Failed to start main window: %v", err)
		os.Exit(1)
	}

	mainWindow.ShowAll()
}

func newSession() {
	dc, err := gui.MakeDeviceChooser(localizer, mainWindow, deviceSelected, logger)
	if err != nil {
		logger.Warnf("Failed to create device chooser: %v", err)
	}

	result := dc.Show()
	logger.Debugf("deviceChooser result: %v", result)
}

func openSession() {
	logger.Debug("OPEN!!!")
}

func saveSession() {
	logger.Debug("SAVE!!!")
}

func saveSessionAs() {
	logger.Debug("SAVE AS!!!")
}

func deviceSelected(device *usb.UsbDevice) {
	logger.Debug("DEVICE SELECTED: %v", device)
}
