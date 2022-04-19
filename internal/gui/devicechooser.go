package gui

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/asaskevich/EventBus"
	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/usb"
)

type DeviceChooser struct {
	window fyne.Window
	bus    *EventBus.Bus
}

func MakeDeviceChooser(a fyne.App, l *localizations.Localizer, bus *EventBus.Bus) *DeviceChooser {

	w := a.NewWindow(l.Get("devicechooser.title"))

	w.Resize(fyne.NewSize(1024, 600))

	lbl := widget.NewLabel(l.Get("devicechooser.label"))

	devices, err := getDevices(l)
	if err != nil {
		return nil
	}

	deviceTable := makeDeviceList(devices, l)

	cancel := widget.NewButton("Cancel", func() {
		w.Hide()
	})
	ok := widget.NewButton("OK", func() {})

	content := container.NewBorder(
		lbl,
		container.NewHBox(layout.NewSpacer(), cancel, ok),
		nil,
		nil,
		deviceTable,
	)

	w.SetContent(content)

	(*bus).Subscribe(NewSessionMsg, func() {
		w.Show()
	})

	return &DeviceChooser{
		window: w,
		bus:    bus,
	}
}

func (dc *DeviceChooser) Show() {
	dc.window.Show()
}

func makeLabel() *widget.Label {
	l := widget.NewLabel("s")
	l.TextStyle.Monospace = true
	return l
}

func makeDeviceList(devices []usb.UsbDevice, l *localizations.Localizer) *widget.List {

	deviceList := widget.NewList(
		func() int {
			return len(devices)
		},
		func() fyne.CanvasObject {
			cnt := container.NewHBox(
				makeLabel(),
				makeLabel(),
				makeLabel(),
			)

			return cnt
		},
		func(idx int, co fyne.CanvasObject) {
			dev := devices[idx]

			(co.(*fyne.Container).Objects[0]).(*widget.Label).SetText(fmt.Sprintf("%s.%s", dev.Bus, dev.Device))
			(co.(*fyne.Container).Objects[1]).(*widget.Label).SetText(fmt.Sprintf("%s:%s", dev.VendorID, dev.ProductID))
			(co.(*fyne.Container).Objects[2]).(*widget.Label).SetText(fmt.Sprintf("%s - %s ", dev.VendorName, dev.ProductName))
		},
	)

	return deviceList
}

func getDevices(l *localizations.Localizer) ([]usb.UsbDevice, error) {
	devices, err := usb.GetDeviceList()

	if err != nil {
		return nil, fmt.Errorf("Error getting device list: %w\n", err)
	}

	sort.Slice(devices, func(i, j int) bool {
		idx := strings.Compare(devices[i].Bus, devices[j].Bus)

		if idx == 0 {
			idx = strings.Compare(devices[i].Device, devices[j].Device)
		}

		return idx <= 0
	})

	return devices, nil
}

// func pad(str string, length int) string {
// 	if len(str) >= length {
// 		return str
// 	}

// 	repeat := length - len(str)
// 	return str + strings.Repeat(" ", repeat)
// }

// func getColWidths(devices []usb.UsbDevice, l *localizations.Localizer) []int {
// 	cols := []int{
// 		len(l.Get("devicechooser.bus")),
// 		len(l.Get("devicechooser.device")),
// 		len(l.Get("devicechooser.vendor-id")),
// 		len(l.Get("devicechooser.product-id")),
// 		len(l.Get("devicechooser.vendor")),
// 		len(l.Get("devicechooser.product")),
// 	}

// 	for _, dev := range devices {
// 		cols[0] = longest(cols[0], dev.Bus)
// 		cols[1] = longest(cols[1], dev.Device)
// 		cols[2] = longest(cols[2], dev.VendorID)
// 		cols[3] = longest(cols[3], dev.ProductID)
// 		cols[4] = longest(cols[4], dev.VendorName)
// 		cols[5] = longest(cols[5], dev.ProductName)
// 	}

// 	return cols
// }

// func longest(curlen int, t string) int {
// 	if len(t) > curlen {
// 		return len(t)
// 	}

// 	return curlen
// }

// func measure(str string) float32 {
// 	style := fyne.TextStyle{
// 		Bold: true,
// 	}
// 	return fyne.MeasureText(str, theme.TextSize(), style).Width
// }
