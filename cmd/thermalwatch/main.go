// thermalwatch prints thermal zone temperatures.
package main

import (
	"fmt"
	"strings"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	zones, err := sys.ChildDirs("/sys/class/thermal")
	sys.Die(err)
	found := false
	for _, zone := range zones {
		if !strings.HasPrefix(zone, "thermal_zone") {
			continue
		}
		base := filepath.Join("/sys/class/thermal", zone)
		temp, err := sys.ReadTrimmed(filepath.Join(base, "temp"))
		if err != nil {
			continue
		}
		found = true
		fmt.Printf("%s %s\n", zone, temp)
	}
	if !found {
		fmt.Println("no thermal zones")
	}
}
