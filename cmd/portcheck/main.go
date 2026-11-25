// portcheck connects to a host and port with a short timeout.
package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	if len(os.Args) != 3 {
		sys.Usage("HOST PORT")
	}
	host, port := os.Args[1], os.Args[2]
	addr := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		fmt.Printf("%s closed: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("%s open\n", addr)
}
