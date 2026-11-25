// neterrors prints interface errors and drops.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/dev")
	sys.Die(err)
	fmt.Printf("%-16s %10s %10s %10s %10s\n", "IFACE", "RX_ERR", "RX_DROP", "TX_ERR", "TX_DROP")
	for _, row := range sys.ParseNetDev(text) {
		fmt.Printf("%-16s %10d %10d %10d %10d\n", row.Name, row.RxErrs, row.RxDrop, row.TxErrs, row.TxDrop)
	}
}
