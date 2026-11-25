// inodewatch prints inode use for each real filesystem mount.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	mounts, err := sys.Mounts()
	sys.Die(err)
	fmt.Printf("%-20s %8s %12s %12s\n", "MOUNT", "USE%", "USED", "TOTAL")
	seen := map[string]bool{}
	for _, m := range mounts {
		if seen[m.Path] || virtual(m.Type) {
			continue
		}
		seen[m.Path] = true
		space, err := sys.Statfs(m.Path)
		if err != nil || space.Files == 0 {
			continue
		}
		used := space.Files - space.FilesFree
		fmt.Printf("%-20s %7.1f%% %12d %12d\n", m.Path, sys.Percent(used, space.Files), used, space.Files)
	}
}

func virtual(kind string) bool {
	switch kind {
	case "proc", "sysfs", "devtmpfs", "devpts", "cgroup", "cgroup2", "pstore", "bpf", "tracefs", "securityfs", "debugfs", "configfs", "fusectl", "mqueue", "hugetlbfs":
		return true
	default:
		return false
	}
}
