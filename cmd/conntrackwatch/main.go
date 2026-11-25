// conntrackwatch prints connection-tracking counters when the kernel exposes them.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	paths := []string{"/proc/sys/net/netfilter/nf_conntrack_count", "/proc/sys/net/netfilter/nf_conntrack_max"}
	for _, path := range paths {
		text, err := sys.ReadTrimmed(path)
		if err != nil {
			fmt.Printf("%s unavailable: %v\n", path, err)
			continue
		}
		fmt.Printf("%s %s\n", path, text)
	}
	text, err := sys.Read("/proc/net/stat/nf_conntrack")
	if err != nil {
		fmt.Printf("conntrack stat unavailable: %v\n", err)
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		fmt.Println(line)
	}
}
