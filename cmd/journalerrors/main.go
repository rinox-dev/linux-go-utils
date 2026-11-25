// journalerrors reports whether a systemd journal exists and scans text logs for errors.
package main

import (
	"fmt"
	"os"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	for _, path := range []string{"/run/log/journal", "/var/log/journal"} {
		if _, err := os.Stat(path); err != nil {
			fmt.Printf("%s absent\n", path)
			continue
		}
		fmt.Printf("%s present\n", path)
	}
	hits, path, err := sys.ScanLines([]string{"/var/log/syslog", "/var/log/messages"}, func(line string) bool {
		lower := strings.ToLower(line)
		return strings.Contains(lower, "error") || strings.Contains(lower, "fail")
	}, 30)
	if err != nil {
		fmt.Printf("text log unavailable: %v\n", err)
		return
	}
	fmt.Printf("source: %s\n", path)
	for _, line := range hits {
		fmt.Println(line)
	}
}
