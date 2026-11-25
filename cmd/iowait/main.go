// iowait prints the I/O wait share of each CPU from /proc/stat.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/stat")
	sys.Die(err)
	fmt.Printf("%-8s %12s %12s %8s\n", "CPU", "IOWAIT", "TOTAL", "PCT")
	for _, cpu := range sys.ParseCPUStat(text) {
		total := cpu.Total()
		fmt.Printf("%-8s %12d %12d %7.1f%%\n", cpu.Name, cpu.IOWait, total, sys.Percent(cpu.IOWait, total))
	}
}
