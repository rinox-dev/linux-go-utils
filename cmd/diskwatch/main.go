// diskwatch prints filesystem use for each mount.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	mounts, err := sys.Mounts()
	sys.Die(err)
	fmt.Printf("%-20s %8s %10s %10s %s\n", "MOUNT", "USE%", "USED", "TOTAL", "TYPE")
	seen := map[string]bool{}
	for _, m := range mounts {
		if seen[m.Path] || skipFS(m.Type) {
			continue
		}
		seen[m.Path] = true
		space, err := sys.Statfs(m.Path)
		if err != nil || space.Total == 0 {
			continue
		}
		fmt.Printf("%-20s %7.1f%% %10s %10s %s\n", m.Path, sys.Percent(space.Used, space.Total), sys.Human(space.Used), sys.Human(space.Total), m.Type)
	}
}

func skipFS(kind string) bool {
	switch kind {
	case "proc", "sysfs", "devtmpfs", "devpts", "cgroup", "cgroup2", "pstore", "bpf", "tracefs", "securityfs", "debugfs", "configfs", "fusectl", "mqueue", "hugetlbfs":
		return true
	default:
		return false
	}
}
