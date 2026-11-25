// loadwatch prints load averages and runnable tasks.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	fmt.Printf("load1: %.2f\nload5: %.2f\nload15: %.2f\nrunning: %d\ntotal: %d\n", load.One, load.Five, load.Fifteen, load.Running, load.Total)
}
