// sudoaudit prints sudo lines from authentication logs.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	hits, path, err := sys.ScanLines([]string{"/var/log/auth.log", "/var/log/secure"}, func(line string) bool {
		return strings.Contains(strings.ToLower(line), "sudo")
	}, 40)
	if err != nil {
		fmt.Printf("log unavailable: %v\n", err)
		return
	}
	fmt.Printf("source: %s\n", path)
	for _, hit := range hits {
		fmt.Println(hit)
	}
	if len(hits) == 0 {
		fmt.Println("no sudo lines")
	}
}
