package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printServiceList(services []Service) {
	if len(services) == 0 {
		fmt.Println("\nThere are no services created yet")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Status", "Endpoint", "Last build"})
		layout := "Jan 2, 2006 at 3:04pm (MST)"
		for _, s := range services {
			v := s.Version.Format(layout)
			table.Append([]string{s.Name, s.Status, s.Endpoint, v})
		}
		table.Render()
	}
}
