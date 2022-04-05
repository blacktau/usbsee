package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/usb"
)

func MakeDeviceChooser(a fyne.App, l *localizations.Localizer) {
	w := a.NewWindow(l.Get("devicechooser.title"))
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(800, 600))

	lbl := widget.NewLabel(l.Get("devicechooser.label"))

	devices, err := usb.GetDeviceList()

	if err != nil {
		fmt.Printf("Error getting device list: %v\n", err)
	}

	tree := widget.NewList(
		func() int {
			return len(devices)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			d := devices[lii]
			co.(*widget.Label).SetText(fmt.Sprintf("%s %s:%s %s %s", d.ID(), d.VendorID, d.ProductID, d.VendorName, d.ProductName))
		},
	)

	content := container.NewBorder(
		lbl,
		nil,
		nil,
		nil,
		tree,
	)

	w.SetContent(content)

	w.Show()
}
