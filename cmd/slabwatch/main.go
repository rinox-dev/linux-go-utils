// slabwatch prints kernel slab cache counters when readable.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/slabinfo")
	if err != nil {
		fmt.Printf("slabinfo unavailable: %v\n", err)
		return
	}
	lines := strings.Split(text, "\n")
	if len(lines) > 2 {
		lines = lines[2:]
	}
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		count++
		if count > 40 {
			break
		}
		fmt.Println(line)
	}
	if count == 0 {
		fmt.Println("no slab entries")
	}
}
