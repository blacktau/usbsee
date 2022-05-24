package usb

import (
	"fmt"
	"strings"

	"github.com/jochenvg/go-udev"
)

type Device struct {
	VendorID    string
	ProductID   string
	VendorName  string
	ProductName string
	Bus         string
	Device      string
}

func (d Device) ID() string {
	return fmt.Sprintf("%s.%s", d.Bus, d.Device)
}

func GetDeviceList() ([]Device, error) {
	return getDeviceListUdev()
}

func getDeviceListUdev() ([]Device, error) {
	u := udev.Udev{}
	en := u.NewEnumerate()
	err := en.AddMatchSubsystem("usb")

	if err != nil {
		return nil, err
	}

	devs, err := en.Devices()

	var result []Device

	for _, d := range devs {
		d.Properties()
		bus := d.PropertyValue("BUSNUM")
		if bus == "" {
			continue
		}

		dev := d.PropertyValue("DEVNUM")
		vid := d.PropertyValue("ID_VENDOR_ID")
		v := d.PropertyValue("ID_VENDOR_FROM_DATABASE")
		pid := d.PropertyValue("ID_MODEL_ID")
		p := getModel(d)

		u := Device{
			VendorID:    vid,
			ProductID:   pid,
			VendorName:  v,
			ProductName: p,
			Bus:         bus,
			Device:      dev,
		}

		result = append(result, u)
	}

	return result, err
}

func getModel(dev *udev.Device) string {
	model := dev.PropertyValue("ID_MODEL_FROM_DATABASE")

	if model == "" {
		model = dev.PropertyValue("ID_MODEL_ENC")
		model = strings.ReplaceAll(model, "\\x20", " ")
	}

	return model
}
