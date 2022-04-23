package gui

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
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

	return appWindow, nil
}

func loadAppIcon() (*gdk.Pixbuf, error) {
	return gdk.PixbufNewFromBytesOnly(logo)
}

func buildMainBody(l *localizations.Localizer) (*gtk.Stack, error) {
	stack, err := gtk.StackNew()
	if err != nil {
		return nil, fmt.Errorf("Could not create packet-stream-stack: %w", err)
	}

	startPage, err := buildStartPage(l)
	if err != nil {
		return nil, fmt.Errorf("error creating start page: %w", err)
	}

	startPage.SetVisible(true)

	stack.AddNamed(startPage, "start-page")

	paned, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, fmt.Errorf("Could not create main stack paned: %w", err)
	}

	paned.SetVisible(true)
	packetStream, err := buildPacketStream()
	if err != nil {
		return nil, fmt.Errorf("Could not create main stack packet stream: %w", err)
	}

	paned.Add1(packetStream)

	stack.AddNamed(paned, "running-page")

	stack.SetVisibleChild(startPage)

	return stack, nil
}

func buildStartPage(l *localizations.Localizer) (*gtk.Box, error) {

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, fmt.Errorf("Could not create start Page: %w", err)
	}

	box.SetVAlign(gtk.ALIGN_CENTER)

	img, err := gtk.ImageNewFromPixbuf(logoBuf)
	if err != nil {
		return nil, fmt.Errorf("Could not create starter image: %w", err)
	}

	box.Add(img)

	label, err := gtk.LabelNew(l.Get("usbsee.start_msg"))

	if err != nil {
		return nil, fmt.Errorf("Could not create starter label: %w", err)
	}

	label.SetHAlign(gtk.ALIGN_CENTER)

	box.PackStart(label, false, true, 0)

	return box, nil
}

func buildPacketStream() (*gtk.ScrolledWindow, error) {
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create packet stream scroller: %w", err)
	}

	list, err := gtk.ListBoxNew()
	if err != nil {
		return nil, fmt.Errorf("Could not create packet list: %w", err)
	}

	scroller.Container.Add(list)

	return scroller, nil
}

func buildHeaderBar(l *localizations.Localizer) (*gtk.HeaderBar, error) {
	header, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, fmt.Errorf("Failed to create application UI: %w", err)
	}

	header.SetTitle(l.Get("usbsee.title"))
	header.SetShowCloseButton(true)

	openButton, err := buildOpenButton(l)
	if err != nil {
		return nil, fmt.Errorf("Failed to create application UI: %w", err)
	}

	header.PackStart(openButton)

	newButton, err := gtk.ButtonNewFromIconName("document-new-symbolic", gtk.ICON_SIZE_BUTTON)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new recording button")
	}

	header.PackStart(newButton)

	saveBtns, err := buildSaveButton(l)

	if err != nil {
		return nil, fmt.Errorf("Failed to create save buttons")
	}

	header.PackEnd(saveBtns)

	return header, nil
}

func buildOpenButton(l *localizations.Localizer) (*gtk.ButtonBox, error) {

	buttonBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, fmt.Errorf("Failed to create open button: %w", err)
	}

	buttonBox.SetLayout(gtk.BUTTONBOX_EXPAND)

	button, err := gtk.ButtonNew()
	if err != nil {
		return nil, fmt.Errorf("Failed to create open button: %w", err)
	}

	button.SetName("open-button")
	button.SetTooltipMarkup(l.Get("usbsee.open_tooltip"))
	button.SetTooltipText(l.Get("usbsee.open_tooltip"))
	button.SetLabel(l.Get("usbsee.open"))
	button.SetUseUnderline(true)
	sc, err := button.GetStyleContext()

	if err != nil {
		return nil, fmt.Errorf("Failed to create open button: %w", err)
	}

	sc.AddClass("linked")

	button.Connect("clicked", func() {
		log.Println("OPEN!!")
	})

	//	buttonBox.Add(button)
	buttonBox.PackStart(button, true, true, 0)

	mru, err := buildMRUButton()
	if err != nil {
		return nil, fmt.Errorf("Failed to build MRU button: %w", err)
	}

	buttonBox.PackStart(mru, false, true, 0)
	buttonBox.SetHomogeneous(false)

	return buttonBox, nil
}

func buildMRUButton() (*gtk.MenuButton, error) {
	menuButton, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, fmt.Errorf("Failed to create MRU button: %w", err)
	}

	menuButton.SetName("open-recent")
	menuButton.SetUsePopover(true)

	return menuButton, nil
}

func buildSaveButton(l *localizations.Localizer) (*gtk.ButtonBox, error) {
	buttonBox, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, fmt.Errorf("Failed to create save button: %w", err)
	}

	buttonBox.SetLayout(gtk.BUTTONBOX_EXPAND)

	saveBtn, err := gtk.ButtonNew()
	if err != nil {
		return nil, fmt.Errorf("Failed to create save button: %w", err)
	}

	saveBtn.SetName("save-button")
	saveBtn.SetTooltipMarkup(l.Get("usbsee.save_tooltip"))
	saveBtn.SetTooltipText(l.Get("usbsee.save_tooltip"))
	saveBtn.SetLabel(l.Get("usbsee.save"))
	saveBtn.SetUseUnderline(true)

	saveBtn.Connect("clicked", func() {
		log.Println("SAVE!!")
	})

	buttonBox.PackStart(saveBtn, true, true, 0)

	saveAs, err := gtk.ButtonNewFromIconName("document-save-as-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, fmt.Errorf("Failed to build Save as button: %w", err)
	}

	saveAs.SetName("save-as-button")
	saveAs.SetTooltipMarkup(l.Get("usbsee.save_as_tooltip"))
	saveAs.SetTooltipText(l.Get("usbsee.save_as_tooltip"))

	buttonBox.PackStart(saveAs, false, true, 0)

	buttonBox.SetHomogeneous(false)

	return buttonBox, nil

}
