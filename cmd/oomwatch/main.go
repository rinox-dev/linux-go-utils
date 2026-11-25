// oomwatch prints the OOM kill counter from /proc/vmstat.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/vmstat")
	sys.Die(err)
	found := false
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && (fields[0] == "oom_kill" || fields[0] == "oom_kill_count" || fields[0] == "pgoom_kill") {
			found = true
			fmt.Printf("%s %s\n", fields[0], fields[1])
		}
	}
	if !found {
		fmt.Println("oom counters unavailable")
	}
}
