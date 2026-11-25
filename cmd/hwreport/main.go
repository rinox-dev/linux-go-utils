// hwreport prints a short hardware summary from procfs and sysfs.
package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	fmt.Printf("memory: %s\n", sys.Human(mem.Total))
	if text, err := sys.Read("/proc/cpuinfo"); err == nil {
		cpus := 0
		for _, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(line, "processor") {
				cpus++
			}
			if strings.HasPrefix(line, "model name") && cpus == 1 {
				_, val, _ := strings.Cut(line, ":")
				fmt.Printf("cpu: %s\n", strings.TrimSpace(val))
			}
		}
		fmt.Printf("cpus: %d\n", cpus)
	}
	zones, err := sys.ChildDirs("/sys/class/thermal")
	if err != nil {
		return
	}
	for _, zone := range zones {
		if !strings.HasPrefix(zone, "thermal_zone") {
			continue
		}
		temp, err := sys.ReadTrimmed(filepath.Join("/sys/class/thermal", zone, "temp"))
		if err != nil {
			continue
		}
		fmt.Printf("%s %s\n", zone, temp)
	}
}
