// powerwatch prints power supply status and capacity.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/power_supply")
	sys.Die(err)
	found := false
	for _, name := range names {
		base := filepath.Join("/sys/class/power_supply", name)
		status, _ := sys.ReadTrimmed(filepath.Join(base, "status"))
		capacity, _ := sys.ReadTrimmed(filepath.Join(base, "capacity"))
		found = true
		fmt.Printf("%s status=%s capacity=%s%%\n", name, status, capacity)
	}
	if !found {
		fmt.Println("no power supplies")
	}
}
