package ebpfusb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"

	bpf "github.com/iovisor/gobpf/bcc"
)

type DirectionFilter uint8

const (
	Incoming DirectionFilter = iota
	Outgoing
	Both
)

const (
	USB_ENDPOINT_XFERTYPE_MASK    = 0x03
	USB_ENDPOINT_XFER_CONTROL     = 0
	USB_ENDPOINT_XFER_ISOCHRONOUS = 1
	USB_ENDPOINT_XFER_BULK        = 2
	USB_ENDPOINT_XFER_INTERRUPT   = 3
	USB_ENDPOINT_XFER_BULK_STREAM = 4

	IN_MAP = 0x0200
)

type UsbEventIntrnl struct {
	Alen          uint64
	Buflen        uint64
	Vendor        uint16
	Product       uint16
	Endpoint      uint8
	TransferFlags uint32
	BmAttributes  uint8
	Buf           [4096]byte
}

type UsbEvent struct {
	Alen          uint64
	Buflen        uint64
	Vendor        uint16
	Product       uint16
	Endpoint      uint8
	BmAttributes  uint8
	TransferFlags uint32
	Buf           [4096]byte
	Direction     string
	TransferType  string
}

const tmplt string = `
#include <linux/usb.h>

struct data_t {
	u64 alen;
	u64 buflen;
	u16 vendor;
	u16 product;
	u8 bmAttributes;
	u8 endpoint;
	u32 transfer_flags;
	u8 buf [4096];
};

BPF_PERF_OUTPUT(events);
BPF_PERCPU_ARRAY(data_struct, struct data_t, 1);

int monitor_usb_hcd_giveback_urb(struct pt_regs *ctx, struct urb *urb) {
	// Perform a VID/PID check if configured to do so
	%s
	// Perform endpoint type filtering if configured to do so
	%s

	int zero = 0;
	struct data_t* data = data_struct.lookup(&zero);
	if (!data)
		return 0;

	struct usb_device *dev = urb->dev;

	data->vendor = dev->descriptor.idVendor;
	data->product = dev->descriptor.idProduct;
	data->alen = urb->actual_length;
	data->transfer_flags = urb->transfer_flags;
	data->buflen = urb->transfer_buffer_length;
	data->endpoint = urb->ep->desc.bEndpointAddress;
	data->bmAttributes = urb->ep->desc.bmAttributes;

	const u8 bmAttr = urb->ep->desc.bmAttributes;

	// bpf_trace_printk("bmAttributes: %%x\n", bmAttr);

	bpf_probe_read_kernel(&data->buf, sizeof(data->buf), urb->transfer_buffer);
	
	events.perf_submit(ctx, data, sizeof(*data));

	return 0;
}
`

type EventHandler func(UsbEvent)

func Start(vendorID, productID *uint16, directionFilter DirectionFilter, handler EventHandler) {
	fmt.Printf("Start(%v,%v,%v,%v)\n", vendorID, productID, directionFilter, handler)
	vndrfltr := calcVendorCheck(vendorID, productID)
	dirfltr := calcDirectionFilter(directionFilter)

	code := fmt.Sprintf(tmplt, vndrfltr, dirfltr)

	mod := bpf.NewModule(code, []string{})
	defer mod.Close()

	probe, err := mod.LoadKprobe("monitor_usb_hcd_giveback_urb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load kprobe kprobe__usb_hcd_giveback_urb: %s\n", err)
		os.Exit(1)
	}

	err = mod.AttachKprobe("__usb_hcd_giveback_urb", probe, -1)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach monitor_usb_hcd_giveback_urb: %s\n", err)
		os.Exit(1)
	}

	eventTbl := bpf.NewTable(mod.TableId("events"), mod)

	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(eventTbl, channel, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	byteOrder := bpf.GetHostByteOrder()

	go func() {
		var event UsbEventIntrnl
		for {
			data := <-channel
			err := binary.Read(bytes.NewBuffer(data), byteOrder, &event)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}

			evt := UsbEvent{
				Alen:          event.Alen,
				Buflen:        event.Buflen,
				Vendor:        event.Vendor,
				Product:       event.Product,
				Endpoint:      event.Endpoint,
				TransferFlags: event.TransferFlags,
				BmAttributes:  event.BmAttributes + 1,
				Buf:           event.Buf,
				Direction:     getEndpointType(event.TransferFlags),
				// I have no idea why in gobpf this is off by one. bpf_trace_printk says its 3 (for INT) but its 2 here.
				TransferType: getTransferType(event.BmAttributes + 1),
			}

			handler(evt)
		}
	}()

	vs := "unspecified"
	if vendorID != nil {
		vs = fmt.Sprintf("0x%x", *vendorID)
	}

	ps := "unspecified"
	if productID != nil {
		ps = fmt.Sprintf("0x%x", *productID)
	}

	fmt.Printf("Starting capture [VID=%s and PID=%s]\n", vs, ps)
	perfMap.Start()
	<-sig
	perfMap.Stop()
}

func calcDirectionFilter(directionFilter DirectionFilter) string {

	if directionFilter == Incoming {
		return "if (~(urb->transfer_flags & 0x0200)) { return 0; }"
	}

	if directionFilter == Outgoing {
		return "if (urb->transfer_flags & 0x0200) { return 0; }"
	}

	return ""
}

func calcVendorCheck(vendorID, productID *uint16) string {
	if vendorID != nil && productID != nil {
		return fmt.Sprintf("if (urb->dev->descriptor.idVendor != %d || urb->dev->descriptor.idProduct != %d) { return 0; }", *vendorID, *productID)
	}

	if vendorID != nil {
		return fmt.Sprintf("if (urb->dev->descriptor.idVendor != %d) {{ return 0; }}", *vendorID)
	}

	if productID != nil {
		return fmt.Sprintf("if (urb->dev->descriptor.idProduct != %d) {{ return 0; }}", *productID)
	}

	return ""
}

func getEndpointType(transferFlags uint32) string {
	if transferFlags&IN_MAP == 0 {
		return "IN"
	}

	return "OUT"
}

func getTransferType(bmAttributes uint8) string {
	masked := USB_ENDPOINT_XFERTYPE_MASK & bmAttributes

	// fmt.Printf("attr: %x %b & mask: %b = %b ? %b\n", bmAttributes, bmAttributes, USB_ENDPOINT_XFERTYPE_MASK, masked, USB_ENDPOINT_XFER_BULK)

	if masked == USB_ENDPOINT_XFER_CONTROL {
		return "CONTROL"
	}

	if masked == USB_ENDPOINT_XFER_ISOCHRONOUS {
		return "ISOC"
	}

	if masked == USB_ENDPOINT_XFER_BULK {
		return "BULK"
	}

	return "INT"
}
