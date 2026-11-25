// mountwatch prints the mount table with read-only and network flags.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	mounts, err := sys.Mounts()
	sys.Die(err)
	fmt.Printf("%-20s %-20s %-8s %s\n", "DEVICE", "PATH", "RO", "TYPE")
	for _, m := range mounts {
		ro := "no"
		if hasOption(m.Options, "ro") {
			ro = "yes"
		}
		fmt.Printf("%-20s %-20s %-8s %s\n", m.Device, m.Path, ro, m.Type)
	}
}

func hasOption(options, want string) bool {
	for _, part := range strings.Split(options, ",") {
		if part == want {
			return true
		}
	}
	return false
}
