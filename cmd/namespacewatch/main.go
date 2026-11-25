// namespacewatch prints namespace links for the current process.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/proc/self/ns")
	sys.Die(err)
	found := false
	for _, name := range names {
		target, err := os.Readlink(filepath.Join("/proc/self/ns", name))
		if err != nil {
			continue
		}
		found = true
		fmt.Printf("%s %s\n", name, target)
	}
	if !found {
		fmt.Println("namespaces unavailable")
	}
}
