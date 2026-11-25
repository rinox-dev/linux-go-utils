// sysinfo prints basic host information from the Linux procfs and uname.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	host, err := os.Hostname()
	if err != nil {
		fail("hostname", err)
	}

	fmt.Printf("Hostname: %s\n", host)
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				fmt.Printf("OS: %s\n", strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`))
				break
			}
		}
	}
	if out, err := exec.Command("uname", "-r", "-m").Output(); err == nil {
		fields := strings.Fields(string(out))
		if len(fields) >= 2 {
			fmt.Printf("Kernel: %s\nArchitecture: %s\n", fields[0], fields[1])
		}
	}
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			fmt.Printf("Load average: %s %s %s\n", fields[0], fields[1], fields[2])
		}
	}
}

func fail(action string, err error) {
	fmt.Fprintf(os.Stderr, "sysinfo: %s: %v\n", action, err)
	os.Exit(1)
}
