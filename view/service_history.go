/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ernestio/ernest-cli/model"
	"github.com/olekukonko/tablewriter"
)

// PrintServiceHistory : Pretty print for service history
func PrintServiceHistory(services []model.Service) {
	if len(services) == 0 {
		fmt.Println("\nThere are no registered builds for this service")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "Version", "User"})
		num := len(services) + 1
		for _, s := range services {
			num = num - 1
			id := strconv.Itoa(num)
			table.Append([]string{id, s.Name, s.Status, s.Version, s.UserName})
		}
		table.Render()
	}
}
