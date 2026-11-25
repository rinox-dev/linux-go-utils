// userprocs groups process counts by user id.
package main

import (
	"fmt"
	"sort"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	counts := map[int]int{}
	for _, proc := range procs {
		counts[proc.UID]++
	}
	keys := make([]int, 0, len(counts))
	for uid := range counts {
		keys = append(keys, uid)
	}
	sort.Ints(keys)
	for _, uid := range keys {
		fmt.Printf("uid=%d processes=%d\n", uid, counts[uid])
	}
}
