// throttlecheck prints CPU thermal throttling counters.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	cpus, err := sys.ChildDirs("/sys/devices/system/cpu")
	sys.Die(err)
	found := false
	for _, cpu := range cpus {
		if len(cpu) < 4 || cpu[:3] != "cpu" || cpu[3] < '0' || cpu[3] > '9' {
			continue
		}
		for _, path := range []string{
			filepath.Join("/sys/devices/system/cpu", cpu, "thermal_throttle", "core_throttle_count"),
			filepath.Join("/sys/devices/system/cpu", cpu, "thermal_throttle", "package_throttle_count"),
		} {
			text, err := sys.ReadTrimmed(path)
			if err != nil {
				continue
			}
			found = true
			fmt.Printf("%s %s %s\n", cpu, filepath.Base(path), text)
		}
	}
	if !found {
		fmt.Println("thermal throttling unavailable")
	}
}
