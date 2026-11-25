// cpufreqwatch prints the current frequency of each CPU.
package main

import (
	"fmt"
	"path/filepath"
	"strconv"

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
		path := filepath.Join("/sys/devices/system/cpu", cpu, "cpufreq", "scaling_cur_freq")
		text, err := sys.ReadTrimmed(path)
		if err != nil {
			continue
		}
		khz, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			fmt.Printf("%s %s kHz\n", cpu, text)
			continue
		}
		found = true
		fmt.Printf("%s %.0f MHz\n", cpu, float64(khz)/1000)
	}
	if !found {
		fmt.Println("cpufreq unavailable")
	}
}
