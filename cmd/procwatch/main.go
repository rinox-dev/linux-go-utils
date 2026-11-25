// procwatch prints process counts by state.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	counts := map[byte]int{}
	for _, proc := range procs {
		counts[proc.State]++
	}
	fmt.Printf("total: %d\n", len(procs))
	for _, state := range []byte{'R', 'S', 'I', 'D', 'Z', 'T'} {
		if counts[state] > 0 {
			fmt.Printf("%c: %d\n", state, counts[state])
		}
	}
}
