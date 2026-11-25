// monitor lists the utilities and runs the matching command binary.
package main

import (
	"fmt"
	"os"

	"linux-go-utils/internal/catalog"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 || (len(args) == 1 && args[0] != "list" && args[0] != "help" && args[0] != "-h" && args[0] != "--help") {
		fmt.Fprintln(os.Stderr, "monitor only lists commands; run each standalone executable directly")
		fmt.Fprintln(os.Stderr, "usage: monitor [list]")
		os.Exit(2)
	}
	for _, c := range catalog.All() {
		fmt.Printf("%-20s %s\n", c.Name, c.Description)
	}
}
