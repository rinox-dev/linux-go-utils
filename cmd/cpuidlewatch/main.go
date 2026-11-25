// cpuidlewatch prints time spent in each CPU idle state.
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
		idle := filepath.Join("/sys/devices/system/cpu", cpu, "cpuidle")
		states, err := sys.ChildDirs(idle)
		if err != nil {
			continue
		}
		for _, state := range states {
			name, err := sys.ReadTrimmed(filepath.Join(idle, state, "name"))
			if err != nil {
				continue
			}
			spent, err := sys.ReadTrimmed(filepath.Join(idle, state, "time"))
			if err != nil {
				continue
			}
			found = true
			fmt.Printf("%s %s %s %s us\n", cpu, state, name, spent)
		}
	}
	if !found {
		fmt.Println("cpuidle unavailable")
	}
}
