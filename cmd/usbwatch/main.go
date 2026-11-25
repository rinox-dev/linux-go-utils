// usbwatch prints USB device names.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/bus/usb/devices")
	sys.Die(err)
	found := false
	for _, name := range names {
		base := filepath.Join("/sys/bus/usb/devices", name)
		vendor, _ := sys.ReadTrimmed(filepath.Join(base, "idVendor"))
		product, _ := sys.ReadTrimmed(filepath.Join(base, "idProduct"))
		if vendor == "" {
			continue
		}
		found = true
		fmt.Printf("%s vendor=%s product=%s\n", name, vendor, product)
	}
	if !found {
		fmt.Println("no USB devices")
	}
}
