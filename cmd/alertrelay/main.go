// alertrelay posts a host snapshot to a webhook URL given on the command line.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"linux-go-utils/internal/sys"
)

func main() {
	if len(os.Args) != 2 {
		sys.Usage("WEBHOOK_URL")
	}
	url := os.Args[1]
	if len(url) < 8 || (url[:7] != "http://" && url[:8] != "https://") {
		fmt.Fprintln(os.Stderr, "alertrelay: webhook must be an http or https URL")
		os.Exit(2)
	}
	load, err := sys.LoadLoadavg()
	sys.Die(err)
	mem, err := sys.LoadMeminfo()
	sys.Die(err)
	host, err := os.Hostname()
	sys.Die(err)
	body, err := json.Marshal(map[string]any{
		"host":                   host,
		"time":                   time.Now().UTC().Format(time.RFC3339),
		"load1":                  load.One,
		"memory_total_bytes":     mem.Total,
		"memory_available_bytes": mem.Available,
	})
	sys.Die(err)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	sys.Die(err)
	defer resp.Body.Close()
	fmt.Printf("status: %s\n", resp.Status)
	if resp.StatusCode >= 300 {
		os.Exit(1)
	}
}
