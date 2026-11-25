// watchdogcheck prints kernel watchdog devices.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/watchdog")
	if err != nil {
		fmt.Printf("watchdog unavailable: %v\n", err)
		return
	}
	found := false
	for _, name := range names {
		base := filepath.Join("/sys/class/watchdog", name)
		identity, _ := sys.ReadTrimmed(filepath.Join(base, "identity"))
		timeout, _ := sys.ReadTrimmed(filepath.Join(base, "timeout"))
		fmt.Printf("%s identity=%s timeout=%s\n", name, identity, timeout)
		found = true
	}
	if _, err := os.Stat("/dev/watchdog"); err == nil {
		fmt.Println("/dev/watchdog present")
		found = true
	}
	if !found {
		fmt.Println("no watchdog")
	}
}
