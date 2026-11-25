// filegrowth lists the largest regular files under a path.
package main

import (
	"fmt"
	"os"

	"linux-go-utils/internal/sys"
)

func main() {
	root := "."
	if len(os.Args) == 2 {
		root = os.Args[1]
	} else if len(os.Args) > 2 {
		sys.Usage("[PATH]")
	}
	total, files, truncated, err := sys.WalkSizes(root, 5000)
	sys.Die(err)
	fmt.Printf("path: %s\ntotal: %s\nfiles: %d\n", root, sys.Human(uint64(total)), len(files))
	if truncated {
		fmt.Println("truncated: yes")
	}
	limit := len(files)
	if limit > 15 {
		limit = 15
	}
	for _, file := range files[:limit] {
		fmt.Printf("%10s %s\n", sys.Human(uint64(file.Size)), file.Path)
	}
}
