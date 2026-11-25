// psfind finds processes whose command line contains the search term.
package main

import (
	"fmt"
	"os"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	if len(os.Args) != 2 {
		sys.Usage("SEARCH_TERM")
	}
	needle := strings.ToLower(os.Args[1])
	procs, err := sys.Processes()
	sys.Die(err)
	found := false
	for _, proc := range procs {
		if strings.Contains(strings.ToLower(proc.Cmdline), needle) {
			fmt.Printf("%-8d %s\n", proc.PID, proc.Cmdline)
			found = true
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "no process matched %q\n", os.Args[1])
		os.Exit(1)
	}
}
