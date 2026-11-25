// failedunits lists systemd services that were started but have no cgroup left.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	entries, err := os.ReadDir("/run/systemd/units")
	if err != nil {
		fmt.Printf("systemd units unavailable: %v\n", err)
		return
	}
	active := map[string]bool{}
	_ = filepath.WalkDir("/sys/fs/cgroup", func(path string, d os.DirEntry, err error) error {
		if err == nil && d.IsDir() && strings.HasSuffix(d.Name(), ".service") {
			active[d.Name()] = true
		}
		return nil
	})
	missing := 0
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "invocation:") || !strings.HasSuffix(name, ".service") {
			continue
		}
		unit := unescape(strings.TrimPrefix(name, "invocation:"))
		if active[unit] {
			continue
		}
		missing++
		fmt.Println(unit)
	}
	if missing == 0 {
		fmt.Println("no invoked service is missing a cgroup")
	}
}

func unescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) && s[i+1] == 'x' {
			n, err := parseHex(s[i+2 : i+4])
			if err == nil {
				b.WriteByte(n)
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func parseHex(s string) (byte, error) {
	var n byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		var d byte
		switch {
		case c >= '0' && c <= '9':
			d = c - '0'
		case c >= 'a' && c <= 'f':
			d = c - 'a' + 10
		case c >= 'A' && c <= 'F':
			d = c - 'A' + 10
		default:
			return 0, fmt.Errorf("bad hex")
		}
		n = n<<4 | d
	}
	return n, nil
}
