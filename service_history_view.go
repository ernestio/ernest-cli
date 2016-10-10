/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printServiceHistory(services []Service) {
	if len(services) == 0 {
		fmt.Println("\nThere are no registered builds for this service")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Build ID", "Status", "Last build"})
		for _, s := range services {
			table.Append([]string{s.Name, s.ID, s.Status, s.Version})
		}
		table.Render()
	}
}
