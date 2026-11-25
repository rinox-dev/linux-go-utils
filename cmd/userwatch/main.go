// userwatch prints local accounts from the passwd database.
package main

import (
	"fmt"
	"os"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/etc/passwd")
	sys.Die(err)
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}
		fmt.Printf("%-16s uid=%s gid=%s shell=%s\n", fields[0], fields[2], fields[3], fields[6])
	}
	_ = os.Args
}
