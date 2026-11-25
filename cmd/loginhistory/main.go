// loginhistory prints recent login records from wtmp.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	records, err := sys.ReadUtmp("/var/log/wtmp")
	sys.Die(err)
	if len(records) > 40 {
		records = records[len(records)-40:]
	}
	for _, rec := range records {
		if rec.User == "" && rec.Type == 8 {
			continue
		}
		fmt.Printf("%-16s %-12s %-16s %s\n", rec.User, rec.Line, rec.Host, rec.When.Format("2006-01-02 15:04"))
	}
	if len(records) == 0 {
		fmt.Println("no login history")
	}
}
