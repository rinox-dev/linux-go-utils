// zombiecheck counts zombie processes.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	n := 0
	for _, proc := range procs {
		if proc.State == 'Z' {
			n++
			fmt.Printf("%d %s\n", proc.PID, proc.Cmdline)
		}
	}
	fmt.Printf("zombies: %d\n", n)
}
