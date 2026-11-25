// tlsdays prints the days until a TLS certificate expires.
package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	host := os.Args[1]
	if len(os.Args) > 2 {
		sys.Usage("HOST[:PORT]")
	}
	if !strings.Contains(host, ":") {
		host += ":443"
	}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", host, &tls.Config{InsecureSkipVerify: true})
	sys.Die(err)
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		fmt.Println("no certificate")
		return
	}
	cert := certs[0]
	days := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
	fmt.Printf("subject: %s\nnotAfter: %s\ndays: %d\n", cert.Subject.CommonName, cert.NotAfter.Format("2006-01-02"), days)
}
