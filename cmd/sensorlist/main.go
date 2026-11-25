// sensorlist prints hardware sensor names from hwmon.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/class/hwmon")
	sys.Die(err)
	for _, name := range names {
		base := filepath.Join("/sys/class/hwmon", name)
		hmName, _ := sys.ReadTrimmed(filepath.Join(base, "name"))
		fmt.Printf("%s %s\n", name, hmName)
	}
}
