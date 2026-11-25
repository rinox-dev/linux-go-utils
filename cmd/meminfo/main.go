// meminfo prints memory and swap counters.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	fmt.Printf("Total: %s\nFree: %s\nAvailable: %s\nBuffers: %s\nCached: %s\nSwapTotal: %s\nSwapFree: %s\n", sys.Human(mem.Total), sys.Human(mem.Free), sys.Human(mem.Available), sys.Human(mem.Buffers), sys.Human(mem.Cached), sys.Human(mem.SwapTotal), sys.Human(mem.SwapFree))
}
