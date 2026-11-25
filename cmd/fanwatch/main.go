// fanwatch prints fan speeds from hwmon.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	hmons, err := sys.ChildDirs("/sys/class/hwmon")
	if err != nil {
		fmt.Printf("hwmon unavailable: %v\n", err)
		return
	}
	found := false
	for _, hm := range hmons {
		base := filepath.Join("/sys/class/hwmon", hm)
		name, _ := sys.ReadTrimmed(filepath.Join(base, "name"))
		for i := 1; i <= 8; i++ {
			text, err := sys.ReadTrimmed(filepath.Join(base, fmt.Sprintf("fan%d_input", i)))
			if err != nil {
				continue
			}
			found = true
			fmt.Printf("%s %s fan%d %s rpm\n", hm, name, i, text)
		}
	}
	if !found {
		fmt.Println("no fan sensors")
	}
}
