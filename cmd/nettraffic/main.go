// nettraffic prints interface byte and packet counters.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/dev")
	sys.Die(err)
	fmt.Printf("%-16s %12s %12s %12s %12s\n", "IFACE", "RX_BYTES", "RX_PKTS", "TX_BYTES", "TX_PKTS")
	for _, row := range sys.ParseNetDev(text) {
		fmt.Printf("%-16s %12d %12d %12d %12d\n", row.Name, row.RxBytes, row.RxPackets, row.TxBytes, row.TxPackets)
	}
}
