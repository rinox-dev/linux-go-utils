// kernelwarn prints warning and error lines from the kernel ring buffer.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	lines, err := sys.ReadKmsg(200)
	if err != nil {
		fmt.Printf("kernel log unavailable: %v\n", err)
		return
	}
	n := 0
	for _, line := range lines {
		if line.Priority > 4 {
			continue
		}
		fmt.Printf("<%d> %s\n", line.Priority, line.Text)
		n++
	}
	if n == 0 {
		fmt.Println("no kernel warnings")
	}
}
