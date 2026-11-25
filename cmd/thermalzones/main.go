// thermalzones prints thermal zone names and temperatures.
package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	zones, err := sys.ChildDirs("/sys/class/thermal")
	sys.Die(err)
	for _, zone := range zones {
		if !strings.HasPrefix(zone, "thermal_zone") {
			continue
		}
		base := filepath.Join("/sys/class/thermal", zone)
		typ, _ := sys.ReadTrimmed(filepath.Join(base, "type"))
		temp, _ := sys.ReadTrimmed(filepath.Join(base, "temp"))
		fmt.Printf("%s type=%s temp=%s\n", zone, typ, temp)
	}
}
