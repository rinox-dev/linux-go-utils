// procage prints the age of processes since start.
package main

import (
	"fmt"
	"sort"
	"time"

	"linux-go-utils/internal/sys"
)

type ageRow struct {
	pid int
	age time.Duration
	cmd string
}

func main() {
	boot, err := sys.BootTime()
	sys.Die(err)
	procs, err := sys.Processes()
	sys.Die(err)
	rows := make([]ageRow, 0, len(procs))
	for _, proc := range procs {
		start := boot.Add(sys.Ticks(proc.StartTicks))
		rows = append(rows, ageRow{pid: proc.PID, age: time.Since(start).Round(time.Second), cmd: proc.Cmdline})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].age > rows[j].age })
	limit := len(rows)
	if limit > 20 {
		limit = 20
	}
	for _, row := range rows[:limit] {
		fmt.Printf("%-8d %-12s %s\n", row.pid, row.age, row.cmd)
	}
}
