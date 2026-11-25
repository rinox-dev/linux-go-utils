// logwatch prints recent error-like lines from a log directory.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	root := "/var/log"
	if len(os.Args) == 2 {
		root = os.Args[1]
	} else if len(os.Args) > 2 {
		sys.Usage("[PATH]")
	}
	var lines []string
	var source string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".log") && !strings.Contains(path, "log") {
			return nil
		}
		text, err := sys.Read(path)
		if err != nil {
			return nil
		}
		for _, line := range strings.Split(text, "\n") {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "error") || strings.Contains(lower, "fail") {
				lines = append(lines, line)
				source = path
			}
		}
		return nil
	})
	if len(lines) == 0 {
		fmt.Println("no matching lines")
		return
	}
	if len(lines) > 40 {
		lines = lines[len(lines)-40:]
	}
	fmt.Printf("source: %s\n", source)
	for _, line := range lines {
		fmt.Println(line)
	}
}
