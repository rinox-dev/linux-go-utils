// serviceprobe connects to one or more host:port endpoints.
package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	if len(os.Args) < 2 {
		sys.Usage("HOST:PORT [HOST:PORT ...]")
	}
	for _, target := range os.Args[1:] {
		conn, err := net.DialTimeout("tcp", target, 3*time.Second)
		if err != nil {
			fmt.Printf("%s closed: %v\n", target, err)
			continue
		}
		conn.Close()
		fmt.Printf("%s open\n", target)
	}
}
