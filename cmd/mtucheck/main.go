// mtucheck prints MTU values for interfaces.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/net")
	sys.Die(err)
	fmt.Printf("%-16s %s\n", "IFACE", "MTU")
	for _, name := range names {
		text, err := sys.ReadTrimmed(filepath.Join("/sys/class/net", name, "mtu"))
		if err != nil {
			continue
		}
		fmt.Printf("%-16s %s\n", name, text)
	}
}
