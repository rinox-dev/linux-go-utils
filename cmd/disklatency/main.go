// disklatency prints block I/O time counters from /proc/diskstats.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/diskstats")
	sys.Die(err)
	fmt.Printf("%-16s %12s %12s %12s %10s\n", "DEVICE", "READ_MS", "WRITE_MS", "IO_MS", "IN_FLIGHT")
	for _, row := range sys.ParseDiskstats(text) {
		if row.ReadMS == 0 && row.WriteMS == 0 && row.IOTicks == 0 {
			continue
		}
		fmt.Printf("%-16s %12d %12d %12d %10d\n", row.Name, row.ReadMS, row.WriteMS, row.IOTicks, row.InFlight)
	}
}
