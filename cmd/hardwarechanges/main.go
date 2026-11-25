// hardwarechanges prints a snapshot of device class counts.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	classes, err := sys.ChildDirs("/sys/class")
	sys.Die(err)
	for _, class := range classes {
		names, err := sys.ChildDirs("/sys/class/" + class)
		if err != nil {
			fmt.Printf("%s unavailable: %v\n", class, err)
			continue
		}
		fmt.Printf("%s %d\n", class, len(names))
	}
}
