// deletedopen lists processes holding file descriptors to deleted files.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	found := false
	for _, proc := range procs {
		_, deleted, err := sys.CountFDs(proc.PID)
		if err != nil || deleted == 0 {
			continue
		}
		found = true
		fmt.Printf("%-8d deleted=%d %s\n", proc.PID, deleted, proc.Cmdline)
	}
	if !found {
		fmt.Println("no deleted open files")
	}
}
