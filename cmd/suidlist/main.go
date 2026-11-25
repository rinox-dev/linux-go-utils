// suidlist scans a path for setuid and setgid files.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"linux-go-utils/internal/sys"
)

func main() {
	root := "/usr/bin"
	if len(os.Args) == 2 {
		root = os.Args[1]
	} else if len(os.Args) > 2 {
		sys.Usage("[PATH]")
	}
	n := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Mode()&(os.ModeSetuid|os.ModeSetgid) == 0 {
			return nil
		}
		n++
		if n <= 200 {
			fmt.Printf("%s %s\n", info.Mode().String(), path)
		}
		return nil
	})
	sys.Die(err)
	fmt.Printf("count: %d\n", n)
}
