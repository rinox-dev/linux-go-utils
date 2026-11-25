// bootlog prints recent kernel ring-buffer lines from /dev/kmsg.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	lines, err := sys.ReadKmsg(80)
	if err != nil {
		fmt.Printf("kernel log unavailable: %v\n", err)
		return
	}
	if len(lines) == 0 {
		fmt.Println("kernel log empty")
		return
	}
	for _, line := range lines {
		fmt.Printf("<%d> %s\n", line.Priority, line.Text)
	}
}
