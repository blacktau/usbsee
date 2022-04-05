package usb

import (
	"fmt"
	"strings"

	"github.com/jochenvg/go-udev"
)

type UsbDevice struct {
	VendorID    string
	ProductID   string
	VendorName  string
	ProductName string
	Bus         string
	Device      string
}

func (d UsbDevice) ID() string {
	return fmt.Sprintf("%s.%s", d.Bus, d.Device)
}

func GetDeviceList() ([]UsbDevice, error) {
	return getDeviceList_Udev()
}

func getDeviceList_Udev() ([]UsbDevice, error) {
	u := udev.Udev{}
	en := u.NewEnumerate()
	en.AddMatchSubsystem("usb")

	devs, err := en.Devices()

	result := []UsbDevice{}

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
		// fmt.Printf("Bus %s Device %s ID: %s:%s %s %s\n", bus, dev, vid, pid, v, p)

		u := UsbDevice{
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
