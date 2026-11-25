// authwatch prints recent authentication lines from local text logs.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	hits, path, err := sys.ScanLines([]string{"/var/log/auth.log", "/var/log/secure"}, interesting, 40)
	if err != nil {
		fmt.Printf("auth log unavailable: %v\n", err)
		return
	}
	fmt.Printf("source: %s\n", path)
	if len(hits) == 0 {
		fmt.Println("no authentication lines matched")
		return
	}
	for _, line := range hits {
		fmt.Println(line)
	}
}

func interesting(line string) bool {
	lower := strings.ToLower(line)
	return strings.Contains(lower, "failed") || strings.Contains(lower, "accepted") || strings.Contains(lower, "authentication")
}
