// routewatch prints the IPv4 route table.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/route")
	sys.Die(err)
	routes := sys.ParseRoutes(text)
	fmt.Printf("%-8s %-16s %-16s %-16s %s\n", "IFACE", "DEST", "GATEWAY", "MASK", "METRIC")
	for _, route := range routes {
		fmt.Printf("%-8s %-16s %-16s %-16s %d\n", route.Iface, route.Destination, route.Gateway, route.Mask, route.Metric)
	}
}
