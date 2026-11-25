// sshwatch prints sshd processes and recent failed SSH lines.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	procs, err := sys.Processes()
	sys.Die(err)
	n := 0
	for _, proc := range procs {
		if strings.Contains(proc.Cmdline, "sshd") {
			fmt.Printf("%d %s\n", proc.PID, proc.Cmdline)
			n++
		}
	}
	fmt.Printf("sshd processes: %d\n", n)
	hits, path, err := sys.ScanLines([]string{"/var/log/auth.log", "/var/log/secure"}, func(line string) bool {
		lower := strings.ToLower(line)
		return strings.Contains(lower, "sshd") && (strings.Contains(lower, "failed") || strings.Contains(lower, "invalid") || strings.Contains(lower, "accepted"))
	}, 25)
	if err != nil {
		fmt.Printf("log unavailable: %v\n", err)
		return
	}
	fmt.Printf("source: %s\n", path)
	for _, hit := range hits {
		fmt.Println(hit)
	}
}
