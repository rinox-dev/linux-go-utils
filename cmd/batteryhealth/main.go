// batteryhealth compares design capacity with the last full charge.
package main

import (
	"fmt"
	"path/filepath"
	"strconv"

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
		design := readUint(filepath.Join(base, "charge_full_design"), filepath.Join(base, "energy_full_design"))
		full := readUint(filepath.Join(base, "charge_full"), filepath.Join(base, "energy_full"))
		fmt.Printf("%s design=%d full=%d", name, design, full)
		if design > 0 {
			fmt.Printf(" health=%.1f%%", sys.Percent(full, design))
		}
		fmt.Println()
	}
	if !found {
		fmt.Println("no battery")
	}
}

func readUint(paths ...string) uint64 {
	for _, path := range paths {
		text, err := sys.ReadTrimmed(path)
		if err != nil {
			continue
		}
		n, err := strconv.ParseUint(text, 10, 64)
		if err == nil {
			return n
		}
	}
	return 0
}
