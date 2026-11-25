// netbaseline prints network counters in a stable key=value form.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/dev")
	sys.Die(err)
	for _, row := range sys.ParseNetDev(text) {
		fmt.Printf("%s.rx_bytes=%d %s.tx_bytes=%d %s.rx_packets=%d %s.tx_packets=%d\n", row.Name, row.RxBytes, row.Name, row.TxBytes, row.Name, row.RxPackets, row.Name, row.TxPackets)
	}
}
