// blockqueue prints in-flight I/O and queue length for block devices.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/diskstats")
	sys.Die(err)
	fmt.Printf("%-16s %10s %12s\n", "DEVICE", "IN_FLIGHT", "NR_REQUESTS")
	for _, row := range sys.ParseDiskstats(text) {
		if len(row.Name) == 0 || row.Name[len(row.Name)-1] >= '0' && row.Name[len(row.Name)-1] <= '9' && !isWhole(row.Name) {
			continue
		}
		nr := "-"
		if v, err := sys.ReadTrimmed(filepath.Join("/sys/block", row.Name, "queue", "nr_requests")); err == nil {
			nr = v
		}
		fmt.Printf("%-16s %10d %12s\n", row.Name, row.InFlight, nr)
	}
}

func isWhole(name string) bool {
	_, err := sys.ReadTrimmed(filepath.Join("/sys/block", name, "size"))
	return err == nil
}
