package gui

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/blacktau/usbsee/internal/logging"
	"github.com/blacktau/usbsee/internal/usb"
)

type DeviceSelected func(device *usb.UsbDevice)

type DeviceChooser struct {
	dialog         *gtk.Window
	selectedDevice *usb.UsbDevice
	logger         logging.Logger
}

var (
	//go:embed assets/device_chooser.glade
	deviceChooserGlade string
)

func makeDeviceChooser(parent gtk.IWindow, logger logging.Logger) (*DeviceChooser, error) {

	dc := &DeviceChooser{
		logger: logger,
	}

	builder, err := gtk.BuilderNewFromString(deviceChooserGlade)

	if err != nil {
		return nil, fmt.Errorf("failed to load device-chooser layout")
	}

	signals := map[string]interface{}{
		"cancel-clicked": func() {
			dc.onCancelClicked()
		},
		"select-clicked": func() {
			dc.onSelectClicked()
		},
		"device-row-selected": func() {
			bObj, err := builder.GetObject("device-chooser-select-button")

			if err != nil {
				logger.Errorf("failed to get select button: %v", err)
				return
			}

			btn, err := isButton(bObj)
			if err != nil {
				logger.Errorf("failed to get select-button not a button: %v", err)
				return
			}

			btn.SetSensitive(true)

			dc.onDeviceSelected()
		},
	}

	builder.ConnectSignals(signals)

	dObj, err := builder.GetObject("device-chooser")

	if err != nil {
		return nil, fmt.Errorf("failed to locate device-chooser in layout: %w", err)
	}

	d, err := isWindow(dObj)

	if err != nil {
		return nil, fmt.Errorf("device-chooser not a dialog?!?: %w", err)
	}

	d.SetTransientFor(parent)

	dc.dialog = d

	devices, err := getDevices()
	if err != nil {
		return nil, fmt.Errorf("failed getDevice list: %w", err)
	}

	dlsObj, err := builder.GetObject("device-list-store")

	if err != nil {
		return nil, fmt.Errorf("could not find device-list-store: %w", err)
	}

	listStore, err := isListStore(dlsObj)

	if err != nil {
		return nil, fmt.Errorf("device-list-store not a ListStore?!?: %w", err)
	}

	for _, dev := range devices {
		iter := listStore.Append()
		err = listStore.Set(iter, []int{0, 1, 2, 3, 4, 5}, []interface{}{dev.Bus, dev.Device, dev.VendorID, dev.ProductID, dev.VendorName, dev.ProductName})
		if err != nil {
			logger.Errorf("could not add device '%s' to list: %v", dev.ID(), err)
		}
	}

	return dc, nil
}

func (dc *DeviceChooser) Show() {
	dc.dialog.ShowNow()
}

func (dc *DeviceChooser) Hide() {
	if dc.dialog != nil {
		dc.dialog.Hide()
	}
}

func (dc *DeviceChooser) onCancelClicked() {
	dc.Hide()
}

func (dc *DeviceChooser) onSelectClicked() {
	dc.logger.Debug("Select Clicked")
}

func (dc *DeviceChooser) onDeviceSelected() {
	dc.logger.Debug("Device Selected!")
}

func getDevices() ([]usb.UsbDevice, error) {
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
