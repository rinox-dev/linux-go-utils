// psfind lists Linux processes whose command line contains a search term.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 || strings.TrimSpace(os.Args[1]) == "" {
		fmt.Fprintln(os.Stderr, "usage: psfind SEARCH_TERM")
		os.Exit(2)
	}
	needle := strings.ToLower(os.Args[1])
	matches := 0
	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "psfind: read /proc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%-8s %s\n", "PID", "COMMAND")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := fmt.Sscan(entry.Name(), new(int)); err != nil {
			continue
		}
		pid := entry.Name()
		data, err := os.ReadFile(filepath.Join("/proc", pid, "cmdline"))
		if err != nil || len(data) == 0 {
			continue
		}
		command := strings.ReplaceAll(string(data), "\x00", " ")
		command = strings.TrimSpace(command)
		if strings.Contains(strings.ToLower(command), needle) {
			fmt.Printf("%-8s %s\n", pid, command)
			matches++
		}
	}
	if matches == 0 {
		fmt.Fprintf(os.Stderr, "psfind: no processes matched %q\n", os.Args[1])
	}
}
