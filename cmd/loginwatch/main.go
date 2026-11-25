// loginwatch prints logged-in sessions from utmp.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	records, err := sys.ReadUtmp("/run/utmp")
	if err != nil {
		records, err = sys.ReadUtmp("/var/run/utmp")
		sys.Die(err)
	}
	n := 0
	for _, rec := range records {
		if rec.Type != 7 { // USER_PROCESS
			continue
		}
		n++
		fmt.Printf("%-16s %-12s %-16s %s\n", rec.User, rec.Line, rec.Host, rec.When.Format("2006-01-02 15:04"))
	}
	if n == 0 {
		fmt.Println("no logged-in sessions")
	}
}
