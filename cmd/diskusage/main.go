// diskusage reports filesystem space for a path (defaults to the current directory).
package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: diskusage [PATH]")
		os.Exit(2)
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		fmt.Fprintf(os.Stderr, "diskusage: %v\n", err)
		os.Exit(1)
	}
	unit := uint64(stat.Bsize)
	total := stat.Blocks * unit
	free := stat.Bavail * unit
	used := total - stat.Bfree*unit
	percent := 0.0
	if total > 0 {
		percent = float64(used) * 100 / float64(total)
	}
	fmt.Printf("Path: %s\nTotal: %s\nUsed:  %s (%.1f%%)\nFree:  %s\n", path, human(total), human(used), percent, human(free))
}

func human(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	value := float64(n)
	units := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	for _, suffix := range units {
		value /= unit
		if value < unit || suffix == "PiB" {
			return fmt.Sprintf("%.1f %s", value, suffix)
		}
	}
	return fmt.Sprintf("%.1f PiB", value)
}
