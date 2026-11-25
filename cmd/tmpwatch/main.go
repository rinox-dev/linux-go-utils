// tmpwatch prints temporary and tmpfs mounts.
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
		if m.Type == "tmpfs" || strings.HasPrefix(m.Path, "/tmp") || strings.HasPrefix(m.Path, "/run") {
			found = true
			fmt.Printf("%s %s %s\n", m.Device, m.Path, m.Type)
		}
	}
	if !found {
		fmt.Println("no temporary mounts")
	}
}
