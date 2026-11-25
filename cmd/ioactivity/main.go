// ioactivity prints block device read and write counters.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/diskstats")
	sys.Die(err)
	fmt.Printf("%-16s %12s %12s %12s %12s\n", "DEVICE", "READS", "READ", "WRITES", "WRITTEN")
	for _, row := range sys.ParseDiskstats(text) {
		if row.Reads == 0 && row.Writes == 0 {
			continue
		}
		fmt.Printf("%-16s %12d %12s %12d %12s\n", row.Name, row.Reads, sys.Human(row.ReadSectors*512), row.Writes, sys.Human(row.WriteSectors*512))
	}
}
