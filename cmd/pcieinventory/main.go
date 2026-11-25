// pcieinventory prints PCI device IDs.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/bus/pci/devices")
	sys.Die(err)
	fmt.Printf("%-16s %-10s %-10s %-10s\n", "SLOT", "VENDOR", "DEVICE", "CLASS")
	for _, name := range names {
		base := filepath.Join("/sys/bus/pci/devices", name)
		vendor, _ := sys.ReadTrimmed(filepath.Join(base, "vendor"))
		device, _ := sys.ReadTrimmed(filepath.Join(base, "device"))
		class, _ := sys.ReadTrimmed(filepath.Join(base, "class"))
		fmt.Printf("%-16s %-10s %-10s %-10s\n", name, vendor, device, class)
	}
}
