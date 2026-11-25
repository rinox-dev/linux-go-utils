// fdcount prints the processes with the most open file descriptors.
package main

import (
	"fmt"
	"sort"

	"linux-go-utils/internal/sys"
)

type fdRow struct {
	pid  int
	open int
	cmd  string
}

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	rows := make([]fdRow, 0, len(procs))
	total := 0
	for _, proc := range procs {
		n, _, err := sys.CountFDs(proc.PID)
		if err != nil {
			continue
		}
		total += n
		rows = append(rows, fdRow{pid: proc.PID, open: n, cmd: proc.Cmdline})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].open > rows[j].open })
	fmt.Printf("open_descriptors: %d\n", total)
	limit := len(rows)
	if limit > 15 {
		limit = 15
	}
	for _, row := range rows[:limit] {
		fmt.Printf("%-8d %6d %s\n", row.pid, row.open, row.cmd)
	}
}
