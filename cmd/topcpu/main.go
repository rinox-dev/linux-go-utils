// topcpu prints processes with the most CPU time.
package main

import (
	"fmt"
	"time"
	"sort"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	sort.Slice(procs, func(i, j int) bool { return procs[i].CPUTime() > procs[j].CPUTime() })
	limit := len(procs)
	if limit > 15 {
		limit = 15
	}
	fmt.Printf("%-8s %10s %s\n", "PID", "CPU", "COMMAND")
	for _, proc := range procs[:limit] {
		fmt.Printf("%-8d %10s %s\n", proc.PID, sys.Ticks(proc.CPUTime()).Round(time.Millisecond), proc.Cmdline)
	}
}
