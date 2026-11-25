// clockwatch prints the clock source and time since boot.
package main

import (
	"fmt"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	fmt.Printf("now: %s\n", time.Now().Format(time.RFC3339))
	up, err := sys.Uptime()
	sys.Die(err)
	fmt.Printf("uptime: %s\n", up.Round(time.Second))
	source, err := sys.ReadTrimmed("/sys/devices/system/clocksource/clocksource0/current_clocksource")
	if err != nil {
		fmt.Printf("clocksource unavailable: %v\n", err)
		return
	}
	fmt.Printf("clocksource: %s\n", source)
}
