package gui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"go.uber.org/zap"

	"github.com/blacktau/usbsee/internal/localizations"
	"github.com/blacktau/usbsee/internal/usb"
)

type DeviceSelected func(device *usb.UsbDevice)

type DeviceChooser struct {
	dialog         *gtk.Dialog
	selectedDevice *usb.UsbDevice
	logger         *zap.SugaredLogger
}

func MakeDeviceChooser(l *localizations.Localizer, parent gtk.IWindow, onDeviceSelected DeviceSelected, logger *zap.SugaredLogger) (*DeviceChooser, error) {

	d, err := gtk.DialogNewWithButtons(
		l.Get("devicechooser.title"),
		parent,
		gtk.DIALOG_DESTROY_WITH_PARENT&gtk.DIALOG_MODAL&gtk.DIALOG_USE_HEADER_BAR,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create dialog: %w", err)
	}

	dc := &DeviceChooser{
		dialog: d,
		logger: logger,
	}

	d.SetStartupID("device-chooser")
	d.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)

	headerBar, err := gtk.HeaderBarNew()
	if err != nil {
		return nil, fmt.Errorf("failed to create headerBar: %w", err)
	}

	headerBar.SetTitle(l.Get("devicechooser.title"))

	selectButton, err := makeDialogButton(l, "devicechooser.select", func () {
		onDeviceSelected(dc.selectedDevice)
		dc.Hide()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create selectButton: %w", err)
	}

	selectButton.SetSensitive(false)

	sc, err := selectButton.GetStyleContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create selectButton: %w", err)
	}

	sc.AddClass("suggested-action")

	headerBar.PackEnd(selectButton)

	cancelButton, err := makeDialogButton(l, "devicechooser.cancel", func() {
		dc.selectedDevice = nil
		dc.Hide()
	})

	if err != nil {
		return nil, err
	}

	headerBar.PackStart(cancelButton)
	headerBar.SetVisible(true)

	d.SetTitlebar(headerBar)

	devices, err := getDevices()
	if err != nil {
		return nil, fmt.Errorf("Failed getDevice list: %w", err)
	}

	deviceList, err := dc.makeDeviceListBox(l, &devices)

	if err != nil {
		return nil, fmt.Errorf("failed to create device list: %w", err)
	}

	deviceList.Connect("row-selected", func(ref *gtk.ListBox, selectedRow *gtk.ListBoxRow) {
		fmt.Println(selectedRow.GetName())
		selectButton.SetSensitive(true)

		devId, err := selectedRow.GetName()
		if err != nil {
			logger.Errorf("failed to get deviceId from selected row: %v", err)
			return
		}

		for _, dev := range devices {
			if dev.ID() == devId {
				var devRef = dev
				dc.selectedDevice = &devRef
			}
		}
	})

	box, err := d.GetContentArea()
	if err != nil {
		return nil, fmt.Errorf("failed to create device list: %w", err)
	}

	box.PackStart(deviceList, true, true, 0)

	return dc, nil
}

func makeDialogButton(l *localizations.Localizer, labelKey string, f interface{}) (*gtk.Button, error) {
	button, err := gtk.ButtonNewWithLabel(l.Get(labelKey))

	if err != nil {
		return nil, fmt.Errorf("failed to create Button '%v': %w", labelKey, err)
	}

	button.SetVisible(true)
	button.SetUseUnderline(true)
	button.Connect("clicked", f)
	return button, nil
}

func (dc *DeviceChooser) Show() *usb.UsbDevice {
	dc.dialog.ShowNow()
	return dc.selectedDevice
}

func (dc *DeviceChooser) Hide() {
	if dc.dialog != nil {
		dc.dialog.Hide()
		dc.dialog.Close()
	}
}

func (dc *DeviceChooser) makeDeviceListBox(l *localizations.Localizer, devices *[]usb.UsbDevice) (*gtk.ListBox, error) {
	lb, err := gtk.ListBoxNew()
	if err != nil {
		return nil, fmt.Errorf("Failed to create list box: %w", err)
	}

	for _, dev := range *devices {
		row, err := dc.makeDeviceListBoxRow(l, dev)
		if err != nil {
			return nil, fmt.Errorf("Could not add device '%s' to list box: %w", dev.Device, err)
		}

		row.SetName(dev.ID())
		lb.Add(row)
	}

	lb.SetVisible(true)

	return lb, nil
}

func (dc *DeviceChooser) makeDeviceListBoxRow(l *localizations.Localizer, device usb.UsbDevice) (*gtk.ListBoxRow, error) {
	row, err := gtk.ListBoxRowNew()

	if err != nil {
		return nil, fmt.Errorf("could not create listboxrow for device %v: %w", device, err)
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		return nil, fmt.Errorf("could not create box for device %v: %w", device, err)
	}

	box.SetVisible(true)
	row.Add(box)

	replacements := &localizations.Replacements{
		"busID":      device.Bus,
		"deviceID":   device.Device,
		"vendorID":   device.VendorID,
		"productID":  device.ProductID,
		"vendorName": device.VendorName,
		"deviceName": device.ProductName,
	}

	_ = addLabel("devicechooser.bus", replacements, l, box)
	_ = addLabel("devicechooser.device", replacements, l, box)
	_ = addLabel("devicechooser.vendor-id-product-id", replacements, l, box)
	_ = addLabel("devicechooser.device-name", replacements, l, box)

	row.SetVisible(true)
	row.SetName(fmt.Sprintf("%v:%v", device.VendorID, device.ProductID))

	return row, nil
}

func addLabel(key string, replacements *localizations.Replacements, l *localizations.Localizer, box *gtk.Box) error {
	lbl, err := gtk.LabelNew("<tt>" + l.Get(key, replacements) + "</tt>")

	if err != nil {
		return fmt.Errorf("could not create label '%s' for device: %w", key, err)
	}

	lbl.SetVisible(true)
	lbl.SetUseMarkup(true)

	box.PackStart(lbl, false, true, 6)

	return nil
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
