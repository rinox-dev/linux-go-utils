// healthcheck exits 1 when load, memory, root disk, or zombies look unhealthy.
package main

import (
	"fmt"
	"os"
	"runtime"

	"linux-go-utils/internal/sys"
)

func main() {
	bad := false
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	limit := float64(runtime.NumCPU()) * 2
	fmt.Printf("load1: %.2f limit: %.2f\n", load.One, limit)
	if load.One > limit {
		bad = true
	}
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	avail := sys.Percent(mem.Available, mem.Total)
	fmt.Printf("memory_available: %.1f%%\n", avail)
	if mem.Total > 0 && avail < 10 {
		bad = true
	}
	space, err := sys.Statfs("/")
	sys.Die(err)
	used := sys.Percent(space.Used, space.Total)
	fmt.Printf("root_used: %.1f%%\n", used)
	if used > 90 {
		bad = true
	}
	procs, err := sys.Processes()
	sys.Die(err)
	zombies := 0
	for _, proc := range procs {
		if proc.State == 'Z' {
			zombies++
		}
	}
	fmt.Printf("zombies: %d\n", zombies)
	if zombies > 0 {
		bad = true
	}
	if bad {
		fmt.Println("status: unhealthy")
		os.Exit(1)
	}
	fmt.Println("status: ok")
}
