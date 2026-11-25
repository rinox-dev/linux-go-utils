// httpcheck requests a URL and prints the status code and duration.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	if len(os.Args) != 2 {
		sys.Usage("URL")
	}
	raw := os.Args[1]
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		fmt.Fprintln(os.Stderr, "httpcheck: URL must be http or https")
		os.Exit(2)
	}
	client := &http.Client{Timeout: 8 * time.Second}
	start := time.Now()
	resp, err := client.Get(raw)
	sys.Die(err)
	defer resp.Body.Close()
	n, _ := io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	fmt.Printf("status: %s\nbytes: %d\nduration: %s\n", resp.Status, n, time.Since(start).Round(time.Millisecond))
	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}
