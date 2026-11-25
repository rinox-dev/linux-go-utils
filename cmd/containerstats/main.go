// containerstats reads CPU and memory files from container cgroups.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	found := 0
	err := filepath.WalkDir("/sys/fs/cgroup", func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || found >= 30 {
			return nil
		}
		base := filepath.Base(path)
		if !strings.Contains(base, "docker-") && !strings.Contains(base, "libpod-") && !strings.Contains(base, "cri-containerd-") {
			return nil
		}
		mem, memErr := sys.ReadTrimmed(filepath.Join(path, "memory.current"))
		cpu, cpuErr := sys.ReadTrimmed(filepath.Join(path, "cpu.stat"))
		fmt.Printf("%s\n", path)
		if memErr != nil {
			fmt.Printf("  memory unavailable: %v\n", memErr)
		} else {
			fmt.Printf("  memory.current: %s\n", mem)
		}
		if cpuErr != nil {
			fmt.Printf("  cpu.stat unavailable: %v\n", cpuErr)
		} else {
			fmt.Printf("  cpu.stat: %s\n", strings.ReplaceAll(cpu, "\n", " "))
		}
		found++
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		sys.Die(err)
	}
	if found == 0 {
		fmt.Println("no container cgroups")
	}
}
