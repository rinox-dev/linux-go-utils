// mountlist prints mounted filesystems from /proc/mounts.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/mounts")
	sys.Die(err)
	mounts := sys.ParseMounts(text)
	fmt.Printf("%-20s %-20s %-10s %s\n", "DEVICE", "PATH", "TYPE", "OPTIONS")
	for _, m := range mounts {
		fmt.Printf("%-20s %-20s %-10s %s\n", m.Device, m.Path, m.Type, m.Options)
	}
	_ = strings.TrimSpace
}
