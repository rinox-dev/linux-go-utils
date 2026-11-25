// netwatch prints interface operational state.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/net")
	sys.Die(err)
	fmt.Printf("%-16s %s\n", "IFACE", "STATE")
	for _, name := range names {
		text, err := sys.ReadTrimmed(filepath.Join("/sys/class/net", name, "operstate"))
		if err != nil {
			text = "unknown"
		}
		fmt.Printf("%-16s %s\n", name, text)
	}
}
