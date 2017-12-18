/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintEnvList : Pretty print for a build list
func PrintEnvList(envs []*emodels.Environment) {
	if len(envs) == 0 {
		fmt.Println("\nThere are no environments created yet")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Project", "Status"})
		for _, e := range envs {
			table.Append([]string{e.Name, e.Project, e.Status})
		}
		table.Render()
	}
}
