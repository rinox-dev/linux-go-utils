// listenchanges prints the current listening TCP and UDP sockets.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	socks, err := sys.TCPSockets()
	if err != nil {
		fmt.Printf("tcp unavailable: %v\n", err)
	} else {
		for _, sock := range socks {
			if sock.State == "LISTEN" {
				fmt.Printf("tcp %s\n", sock.Local)
			}
		}
	}
	for _, path := range []string{"/proc/net/udp", "/proc/net/udp6"} {
		text, err := sys.Read(path)
		if err != nil {
			continue
		}
		v6 := strings.HasSuffix(path, "6")
		for _, line := range strings.Split(text, "\n")[1:] {
			f := strings.Fields(line)
			if len(f) < 4 || f[3] != "07" {
				continue
			}
			fmt.Printf("udp %s\n", f[1])
			_ = v6
		}
	}
}
