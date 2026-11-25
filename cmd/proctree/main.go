// proctree prints process parent identifiers.
package main

import (
	"fmt"
	"sort"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	byPID := map[int]sys.Process{}
	for _, proc := range procs {
		byPID[proc.PID] = proc
	}
	rows := append([]sys.Process(nil), procs...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].PID < rows[j].PID })
	for _, proc := range rows {
		parent := "-"
		if p, ok := byPID[proc.PPID]; ok {
			parent = p.Comm
		}
		fmt.Printf("%-8d ppid=%-8d %-12s %s\n", proc.PID, proc.PPID, parent, proc.Cmdline)
	}
}
