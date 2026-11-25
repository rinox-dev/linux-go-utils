// bandwidth samples /proc/net/dev for one second and prints the byte rate.
package main

import (
	"fmt"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	first, err := readDev()
	sys.Die(err)
	time.Sleep(time.Second)
	second, err := readDev()
	sys.Die(err)
	fmt.Printf("%-12s %12s %12s\n", "IFACE", "RX/s", "TX/s")
	for name, b := range second {
		a, ok := first[name]
		if !ok {
			continue
		}
		fmt.Printf("%-12s %12s %12s\n", name, sys.Human(delta(b.RxBytes, a.RxBytes)), sys.Human(delta(b.TxBytes, a.TxBytes)))
	}
}

func readDev() (map[string]sys.NetDev, error) {
	text, err := sys.Read("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	out := map[string]sys.NetDev{}
	for _, row := range sys.ParseNetDev(text) {
		out[row.Name] = row
	}
	return out, nil
}

func delta(now, then uint64) uint64 {
	if now < then {
		return 0
	}
	return now - then
}
