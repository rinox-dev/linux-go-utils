// topmem prints processes with the most resident memory.
package main

import (
	"fmt"
	"sort"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	sort.Slice(procs, func(i, j int) bool { return procs[i].VmRSS > procs[j].VmRSS })
	limit := len(procs)
	if limit > 15 {
		limit = 15
	}
	fmt.Printf("%-8s %12s %s\n", "PID", "RSS", "COMMAND")
	for _, proc := range procs[:limit] {
		fmt.Printf("%-8d %12s %s\n", proc.PID, sys.Human(proc.VmRSS), proc.Cmdline)
	}
}
