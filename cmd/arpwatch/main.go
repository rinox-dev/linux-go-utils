// arpwatch prints the IPv4 neighbor cache from /proc/net/arp.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/arp")
	sys.Die(err)
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) == 0 {
		fmt.Println("no neighbors")
		return
	}
	fmt.Printf("%-16s %-18s %-8s %s\n", "IP", "HWADDR", "FLAGS", "DEVICE")
	for _, line := range lines[1:] {
		f := strings.Fields(line)
		if len(f) < 6 {
			continue
		}
		fmt.Printf("%-16s %-18s %-8s %s\n", f[0], f[3], f[2], f[5])
	}
}
