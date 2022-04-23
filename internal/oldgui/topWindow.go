package oldgui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/asaskevich/EventBus"
	"github.com/blacktau/usbsee/internal/localizations"
)

type TopWindow struct {
	window    fyne.Window
	localizer *localizations.Localizer
	bus       *EventBus.Bus
}

const NewSessionMsg = "new-session"
const SaveSessionMsg = "save-session"
const OpenSessionMsg = "open-session"

func MakeTopWindow(a fyne.App, l *localizations.Localizer, bus *EventBus.Bus) *TopWindow {

	topWindow := a.NewWindow("Usbsee")
	tw := &TopWindow{
		window:    topWindow,
		localizer: l,
		bus:       bus,
	}

	topWindow.Resize(fyne.NewSize(1024.0, 768.0))
	topWindow.SetMainMenu(MakeMainMenu(&a, l))
	topWindow.SetMaster()

	// MakeDeviceChooser(a, l)
	topWindow.SetContent(
		container.NewBorder(
			tw.makeTopBar(),
			nil,
			nil,
			nil,
		))

	return tw
}

func (tw *TopWindow) ShowAndRun() {
	tw.window.ShowAndRun()
}

func (tw *TopWindow) makeToolBar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			(*tw.bus).Publish(NewSessionMsg)
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			(*tw.bus).Publish(SaveSessionMsg)
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			(*tw.bus).Publish(OpenSessionMsg)
		}),
	)
}

func (tw *TopWindow) makeTopBar() *fyne.Container {
	top := tw.makeToolBar()
	return container.NewVBox(top, widget.NewSeparator())
}
