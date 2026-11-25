// cgroupwatch prints the cgroup version and controllers.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	if text, err := sys.ReadTrimmed("/sys/fs/cgroup/cgroup.controllers"); err == nil {
		fmt.Printf("version: 2\ncontrollers: %s\n", text)
	} else {
		fmt.Println("version: 1 or unavailable")
		fmt.Printf("controllers file: %v\n", err)
	}
	text, err := sys.Read("/proc/self/cgroup")
	if err != nil {
		fmt.Printf("self cgroup unavailable: %v\n", err)
		return
	}
	fmt.Println("self:")
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		fmt.Println(line)
	}
}
