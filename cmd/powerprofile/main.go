// powerprofile prints the ACPI platform power profile.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.ReadTrimmed("/sys/firmware/acpi/platform_profile")
	if err != nil {
		fmt.Printf("power profile unavailable: %v\n", err)
		return
	}
	fmt.Printf("profile: %s\n", text)
}
