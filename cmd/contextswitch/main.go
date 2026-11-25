// contextswitch prints the context-switch counter from /proc/stat.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	n, err := sys.StatField("ctxt")
	sys.Die(err)
	fmt.Printf("context_switches: %d\n", n)
}
