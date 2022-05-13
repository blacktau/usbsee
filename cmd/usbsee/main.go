package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/jessevdk/go-flags"

	"github.com/blacktau/usbsee/ebpfusb"
)

var (
	eventNumber int64 = 0
)

type Options struct {
	VendorID        *string `short:"v" long:"vendorid" description:"the vendor-id in hex to filter for"`
	ProductID       *string `short:"p" long:"productid" description:"the product-id in hex to filter for"`
	Truncate        bool    `short:"t" long:"truncate" description:"trim hex packets to their actual length"`
	FilterDirection *string `short:"d" long:"direction" description:"filter to input or output only. valid values in or out" choice:"in" choice:"out"`
}

func main() {

	euid := os.Geteuid()

	if euid != 0 {
		fmt.Fprintln(os.Stderr, "This Program needs to be run as root.")
		os.Exit(1)
	}

	var opts = &Options{}

	parser := flags.NewParser(opts, flags.Default)
	parser.Command.Name = "usbsee"
	parser.Usage = "- A cli tool for monitoring usb traffic that allows filtering using eBPF"

	_, err := parser.Parse()

	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	vID := hexToUintPtr(opts.VendorID)
	pID := hexToUintPtr(opts.ProductID)

	direction := ebpfusb.Both

	if opts.FilterDirection != nil && *opts.FilterDirection == "in" {
		direction = ebpfusb.Incoming
	}

	if opts.FilterDirection != nil && *opts.FilterDirection == "out" {
		direction = ebpfusb.Outgoing
	}

	monitor := ebpfusb.MakeUsbMonitor(vID, pID, direction, printEvent(opts.Truncate))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	vs := "unspecified"
	if vID != nil {
		vs = fmt.Sprintf("0x%x", *vID)
	}

	ps := "unspecified"
	if pID != nil {
		ps = fmt.Sprintf("0x%x", *pID)
	}

	fmt.Fprintf(os.Stdout, "Initialising monitoring on [VID=%s and PID=%s]\n", vs, ps)

	err = monitor.Init()
	defer monitor.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initalise usb monitor: %v", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Starting monitoring on [VID=%s and PID=%s]\n", vs, ps)

	err = monitor.Start()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start monitor: %v", err)
		os.Exit(1)
	}

	<-sig
}

func hexToUintPtr(src *string) *uint16 {
	if src == nil {
		return nil
	}

	s := *src

	val, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse `%s` as hex\n", *src)
		os.Exit(1)
	}

	cnv := uint16(val)
	return &cnv
}

func printEvent(truncate bool) ebpfusb.EventHandler {
	return func(event ebpfusb.UsbEvent) {
		fmt.Fprintf(os.Stdout,
			"%d: %04x:%04x [0x%02x %s] (%s) actual length = %d, buffer length = %d\n",
			eventNumber,
			event.Vendor,
			event.Product,
			event.Endpoint,
			event.Direction,
			event.TransferType,
			event.Alen,
			event.Buflen,
		)

		len := event.Buflen

		if truncate {
			len = event.Alen
		}

		fmt.Fprintln(os.Stdout, hex.Dump(event.Buf[0:len]))
		eventNumber += 1
	}
}
