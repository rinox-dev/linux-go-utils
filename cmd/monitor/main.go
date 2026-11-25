// monitor lists the utilities and runs the matching command binary.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"linux-go-utils/internal/catalog"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "list" || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		for _, c := range catalog.All() {
			fmt.Printf("%-20s %s\n", c.Name, c.Description)
		}
		return
	}
	name := args[0]
	if !catalog.Has(name) {
		fmt.Fprintf(os.Stderr, "unknown utility %q; run monitor list\n", name)
		os.Exit(2)
	}
	bin := filepath.Join(filepath.Dir(os.Args[0]), name)
	cmd := exec.Command(bin, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "monitor: %v\n", err)
		os.Exit(1)
	}
}
