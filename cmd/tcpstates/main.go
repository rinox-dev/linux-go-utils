// tcpstates groups TCP sockets by kernel state.
package main

import (
	"fmt"
	"sort"

	"linux-go-utils/internal/sys"
)

func main() {
	socks, err := sys.TCPSockets()
	sys.Die(err)
	counts := map[string]int{}
	for _, sock := range socks {
		counts[sock.State]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%-16s %d\n", key, counts[key])
	}
	fmt.Printf("total: %d\n", len(socks))
}
