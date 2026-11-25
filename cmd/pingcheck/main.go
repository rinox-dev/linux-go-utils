// pingcheck sends an ICMP echo to a host.
package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	host := os.Args[1]
	if len(os.Args) > 2 {
		sys.Usage("HOST")
	}
	if host == "" {
		sys.Usage("HOST")
	}
	conn, err := net.Dial("ip4:icmp", host)
	if err != nil {
		fmt.Printf("ping failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
		sys.Die(err)
	}
	packet := make([]byte, 8)
	packet[0] = 8 // Echo Request
	packet[1] = 0
	binary.BigEndian.PutUint16(packet[4:6], 1) // id
	binary.BigEndian.PutUint16(packet[6:8], 1) // seq
	checksum := icmpChecksum(packet)
	binary.BigEndian.PutUint16(packet[2:4], checksum)
	start := time.Now()
	if _, err := conn.Write(packet); err != nil {
		fmt.Printf("ping failed: %v\n", err)
		os.Exit(1)
	}
	buf := make([]byte, 1500)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("reply from %s: %d bytes time=%s\n", host, n, time.Since(start).Round(time.Millisecond))
}

func icmpChecksum(b []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(b); i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	if len(b)%2 == 1 {
		sum += uint32(b[len(b)-1]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	return ^uint16(sum)
}
