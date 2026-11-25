// promtext prints selected host metrics in Prometheus text format.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	up, err := sys.Uptime()
	sys.Die(err)
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	fmt.Printf("linux_uptime_seconds %.2f\n", up.Seconds())
	fmt.Printf("linux_load1 %.2f\nlinux_load5 %.2f\nlinux_load15 %.2f\n", load.One, load.Five, load.Fifteen)
	fmt.Printf("linux_memory_total_bytes %d\nlinux_memory_available_bytes %d\n", mem.Total, mem.Available)
}
