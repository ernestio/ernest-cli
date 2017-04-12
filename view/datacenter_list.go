/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ernestio/ernest-cli/model"
	"github.com/olekukonko/tablewriter"
)

// PrintDatacenterList : Pretty print for a datacenter list
func PrintDatacenterList(datacenters []model.Datacenter) {
	if len(datacenters) == 0 {
		fmt.Println("There are no datacenters created yet.")
		return
	}

	var aws []model.Datacenter
	var vcloud []model.Datacenter
	var azure []model.Datacenter

	for _, d := range datacenters {
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
		fmt.Println("AWS Datacenters")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Group", "Type", "Region", "Url"})
		for _, d := range aws {
			id := strconv.Itoa(d.ID)
			table.Append([]string{id, d.Name, d.GroupName, d.Type, d.Region, d.VseURL})
		}
		table.Render()
	}

	if len(vcloud) > 0 {
		fmt.Println("")
		fmt.Println("VCloud Datacenters")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Group", "Type", "Url", "External Network", "Org"})
		for _, d := range vcloud {
			id := strconv.Itoa(d.ID)
			parts := strings.Split(d.Username, "@")
			org := ""
			if len(parts) == 2 {
				org = parts[1]
			}
			table.Append([]string{id, d.Name, d.GroupName, d.Type, d.VCloudURL, d.ExternalNetwork, org})
		}
		table.Render()
	}

	if len(azure) > 0 {
		fmt.Println("")
		fmt.Println("Azure Datacenters")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Group", "Type", "Region"})
		for _, d := range azure {
			id := strconv.Itoa(d.ID)
			table.Append([]string{id, d.Name, d.GroupName, d.Type, d.Region})
		}
		table.Render()
	}
}
