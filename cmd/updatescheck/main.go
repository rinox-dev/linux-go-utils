// updatescheck reports the size of the installed package database without calling package managers.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	for _, item := range []struct {
		name, path string
	}{
		{"pacman", "/var/lib/pacman/local"},
		{"dpkg", "/var/lib/dpkg/status"},
		{"rpm", "/var/lib/rpm"},
	} {
		info, err := os.Stat(item.path)
		if err != nil {
			fmt.Printf("%s unavailable: %v\n", item.name, err)
			continue
		}
		if info.IsDir() {
			entries, err := os.ReadDir(item.path)
			if err != nil {
				fmt.Printf("%s unavailable: %v\n", item.name, err)
				continue
			}
			fmt.Printf("%s packages: %d\n", item.name, len(entries))
		} else {
			fmt.Printf("%s database: %s\n", item.name, sys.Human(uint64(info.Size())))
		}
		_ = filepath.Base
	}
}
