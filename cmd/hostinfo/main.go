// hostinfo prints host, distribution, kernel, and architecture from procfs.
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	host, err := os.Hostname()
	sys.Die(err)
	fmt.Printf("Hostname: %s\n", host)
	if text, err := sys.Read("/etc/os-release"); err == nil {
		for _, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				fmt.Printf("OS: %s\n", strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`))
			}
		}
	}
	if rel, err := sys.ReadTrimmed("/proc/sys/kernel/osrelease"); err == nil {
		fmt.Printf("Kernel: %s\n", rel)
	}
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	if text, err := sys.Read("/proc/cpuinfo"); err == nil {
		model := ""
		cpus := 0
		for _, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Model") {
				_, val, ok := strings.Cut(line, ":")
				if ok && model == "" {
					model = strings.TrimSpace(val)
				}
			}
			if strings.HasPrefix(line, "processor") {
				cpus++
			}
		}
		if model != "" {
			fmt.Printf("CPU: %s\n", model)
		}
		if cpus > 0 {
			fmt.Printf("CPUs: %d\n", cpus)
		}
	}
}
