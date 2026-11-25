// uptime prints system uptime and load average.
package main

import (
	"fmt"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	up, err := sys.Uptime()
	sys.Die(err)
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	fmt.Printf("Uptime: %s\nLoad average: %.2f %.2f %.2f\n", up.Round(time.Second), load.One, load.Five, load.Fifteen)
}
