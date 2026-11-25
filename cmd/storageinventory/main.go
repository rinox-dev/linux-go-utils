// storageinventory prints block devices and sizes.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/block")
	sys.Die(err)
	fmt.Printf("%-16s %12s %s\n", "DEVICE", "SECTORS", "MODEL")
	for _, name := range names {
		base := filepath.Join("/sys/block", name)
		size, err := sys.ReadTrimmed(filepath.Join(base, "size"))
		if err != nil {
			continue
		}
		model, _ := sys.ReadTrimmed(filepath.Join(base, "device", "model"))
		fmt.Printf("%-16s %12s %s\n", name, size, model)
	}
}
