// metricdump prints selected host metrics as JSON.
package main

import (
	"encoding/json"
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	up, err := sys.Uptime()
	sys.Die(err)
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	out := map[string]any{
		"uptime_seconds":         up.Seconds(),
		"load1":                  load.One,
		"load5":                  load.Five,
		"load15":                 load.Fifteen,
		"memory_total_bytes":     mem.Total,
		"memory_available_bytes": mem.Available,
	}
	b, err := json.Marshal(out)
	sys.Die(err)
	fmt.Println(string(b))
}
