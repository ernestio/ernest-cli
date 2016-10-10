/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

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
		for _, s := range services {
			table.Append([]string{s.Name, s.Status, s.Endpoint, s.Version})
		}
		table.Render()
	}
}
