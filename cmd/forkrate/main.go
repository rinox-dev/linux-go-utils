// forkrate prints the process-creation counter from /proc/stat.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	n, err := sys.StatField("processes")
	sys.Die(err)
	up, err := sys.Uptime()
	sys.Die(err)
	fmt.Printf("forks_since_boot: %d\n", n)
	sec := up.Seconds()
	if sec > 0 {
		fmt.Printf("forks_per_second: %.2f\n", float64(n)/sec)
	}
}
