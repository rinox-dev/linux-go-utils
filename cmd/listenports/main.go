// listenports prints listening TCP sockets from the kernel socket tables.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	socks, err := sys.TCPSockets()
	sys.Die(err)
	fmt.Printf("%-24s %s\n", "LOCAL", "STATE")
	n := 0
	for _, sock := range socks {
		if sock.State != "LISTEN" {
			continue
		}
		fmt.Printf("%-24s %s\n", sock.Local, sock.State)
		n++
	}
	if n == 0 {
		fmt.Println("no listeners")
	}
}
