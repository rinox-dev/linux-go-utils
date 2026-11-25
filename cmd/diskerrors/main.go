// diskerrors prints driver error counters exposed under /sys/block.
package main

import (
	"fmt"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	names, err := sys.ChildDirs("/sys/block")
	sys.Die(err)
	found := false
	for _, name := range names {
		base := filepath.Join("/sys/block", name)
		for _, key := range []string{"device/ioerr_cnt", "device/iodone_cnt", "device/iorequest_cnt"} {
			text, err := sys.ReadTrimmed(filepath.Join(base, key))
			if err != nil {
				continue
			}
			found = true
			fmt.Printf("%s %s %s\n", name, filepath.Base(key), text)
		}
	}
	if !found {
		fmt.Println("no disk error counters")
	}
}
