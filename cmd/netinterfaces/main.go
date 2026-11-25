// netinterfaces prints interface addresses and operational state.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/net")
	sys.Die(err)
	fmt.Printf("%-16s %-16s %-12s %s\n", "IFACE", "OPERSTATE", "MTU", "ADDR")
	for _, name := range names {
		base := filepath.Join("/sys/class/net", name)
		oper, _ := sys.ReadTrimmed(filepath.Join(base, "operstate"))
		mtu, _ := sys.ReadTrimmed(filepath.Join(base, "mtu"))
		addr, _ := os.ReadFile(filepath.Join(base, "address"))
		fmt.Printf("%-16s %-16s %-12s %s\n", name, oper, mtu, strings.TrimSpace(string(addr)))
	}
}
