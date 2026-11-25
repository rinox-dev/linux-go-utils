// permcheck prints modes of critical account files.
package main

import (
	"fmt"
	"os"

)

func main() {
	for _, path := range []string{"/etc/passwd", "/etc/shadow", "/etc/sudoers"} {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("%s unavailable: %v\n", path, err)
			continue
		}
		fmt.Printf("%s %s\n", path, info.Mode().String())
	}
}
