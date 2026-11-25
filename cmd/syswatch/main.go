// syswatch prints a compact host and load snapshot.
package main

import (
	"fmt"
	"os"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	up, err := sys.Uptime()
	sys.Die(err)
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	fmt.Printf("host: %s\nuptime: %s\nload1: %.2f\nmemory_available: %s\n", host, up.Round(time.Second), load.One, sys.Human(mem.Available))
}
