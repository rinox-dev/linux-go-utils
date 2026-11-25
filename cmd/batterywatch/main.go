// batterywatch prints battery status and capacity from sysfs.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/power_supply")
	if err != nil {
		fmt.Printf("power supply unavailable: %v\n", err)
		return
	}
	found := false
	for _, name := range names {
		base := filepath.Join("/sys/class/power_supply", name)
		kind, err := sys.ReadTrimmed(filepath.Join(base, "type"))
		if err != nil || kind != "Battery" {
			continue
		}
		found = true
		status, _ := sys.ReadTrimmed(filepath.Join(base, "status"))
		capText, _ := sys.ReadTrimmed(filepath.Join(base, "capacity"))
		fmt.Printf("%s status=%s capacity=%s%%\n", name, status, capText)
	}
	if !found {
		fmt.Println("no battery")
	}
}
