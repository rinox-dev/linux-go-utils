// certwatch reports days left on PEM certificates under a directory.
package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	root := "/etc/ssl/certs"
	if len(os.Args) == 2 {
		root = os.Args[1]
	} else if len(os.Args) > 2 {
		sys.Usage("[DIR]")
	}
	now := time.Now()
	found := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if found >= 200 {
			return filepath.SkipAll
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		rest := b
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			if block.Type != "CERTIFICATE" {
				continue
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}
			days := int(cert.NotAfter.Sub(now).Hours() / 24)
			fmt.Printf("%s %s expires=%s days=%d\n", path, cert.Subject.CommonName, cert.NotAfter.Format("2006-01-02"), days)
			found++
		}
		return nil
	})
	sys.Die(err)
	if found == 0 {
		fmt.Println("no certificates")
	}
}
