/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintProjectList : Pretty print for a project list
func PrintProjectList(projects []*emodels.Project) {
	if len(projects) == 0 {
		fmt.Println("There are no projects created yet.")
		return
	}

	var aws []*emodels.Project
	var vcloud []*emodels.Project
	var azure []*emodels.Project

	for _, d := range projects {
		switch d.Type {
		case "aws", "aws-fake":
			aws = append(aws, d)
		case "vcloud", "vcloud-fake":
			vcloud = append(vcloud, d)
		case "azure", "azure-fake":
			azure = append(azure, d)
		}
	}

	if len(aws) > 0 {
		fmt.Println("")
		fmt.Println("AWS Projects")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Region"})
		for _, d := range aws {
			id := strconv.Itoa(d.ID)
			region, _ := d.Credentials["region"].(string)
			table.Append([]string{id, d.Name, d.Type, region})
		}
		table.Render()
	}

	if len(vcloud) > 0 {
		fmt.Println("")
		fmt.Println("VCloud Projects")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type", "Vdc", "Org", "Url"})
		for _, d := range vcloud {
			id := strconv.Itoa(d.ID)
			vcloudURL, _ := d.Credentials["vcloud_url"].(string)
			vdc, _ := d.Credentials["vdc"].(string)
			username, _ := d.Credentials["username"].(string)
			parts := strings.Split(username, "@")
			org := ""
			if len(parts) == 2 {
				org = parts[1]
			}
			table.Append([]string{id, d.Name, d.Type, vdc, org, vcloudURL})
		}
		table.Render()
	}

	if len(azure) > 0 {
		fmt.Println("")
		fmt.Println("Azure Projects")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Type"})
		for _, d := range azure {
			id := strconv.Itoa(d.ID)
			table.Append([]string{id, d.Name, d.Type})
		}
		table.Render()
	}
}
