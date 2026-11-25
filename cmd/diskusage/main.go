// diskusage reports filesystem space for a path. It uses statfs, not df.
package main

import (
	"fmt"
	"os"

	"linux-go-utils/internal/sys"
)

func main() {
	path := "."
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		sys.Usage("[PATH]")
	}
	space, err := sys.Statfs(path)
	sys.Die(err)
	fmt.Printf("Path: %s\nTotal: %s\nUsed: %s (%.1f%%)\nAvailable: %s\n", path, sys.Human(space.Total), sys.Human(space.Used), sys.Percent(space.Used, space.Total), sys.Human(space.Avail))
}
