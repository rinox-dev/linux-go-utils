// swapwatch prints swap use and paging counters.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	fmt.Printf("SwapTotal: %s\nSwapFree: %s\nSwapUsed: %s\n", sys.Human(mem.SwapTotal), sys.Human(mem.SwapFree), sys.Human(mem.SwapTotal-mem.SwapFree))
	text, err := sys.Read("/proc/vmstat")
	sys.Die(err)
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && (fields[0] == "pswpin" || fields[0] == "pswpout") {
			fmt.Printf("%s %s\n", fields[0], fields[1])
		}
	}
}
