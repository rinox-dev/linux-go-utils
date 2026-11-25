// threadcount prints thread counts by process.
package main

import (
	"fmt"
	"sort"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	type row struct {
		pid int
		n   int
		cmd string
	}
	rows := make([]row, 0, len(procs))
	total := 0
	for _, proc := range procs {
		total += proc.Threads
		rows = append(rows, row{pid: proc.PID, n: proc.Threads, cmd: proc.Cmdline})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].n > rows[j].n })
	fmt.Printf("processes: %d threads: %d\n", len(procs), total)
	limit := len(rows)
	if limit > 15 {
		limit = 15
	}
	for _, row := range rows[:limit] {
		fmt.Printf("%-8d %6d %s\n", row.pid, row.n, row.cmd)
	}
}
