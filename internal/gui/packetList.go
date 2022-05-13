package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func makePacketList() {
	packetListStore, err := gtk.ListStoreNew(glib.TYPE_INT64)

	if err != nil {
		fmt.Printf("Error %v", err)

	}
	fmt.Println(packetListStore)
}
