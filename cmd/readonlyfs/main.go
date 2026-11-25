// readonlyfs prints filesystems mounted read-only.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	mounts, err := sys.Mounts()
	sys.Die(err)
	found := false
	for _, m := range mounts {
		if hasOption(m.Options, "ro") {
			found = true
			fmt.Printf("%s %s\n", m.Device, m.Path)
		}
	}
	if !found {
		fmt.Println("no read-only filesystems")
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
