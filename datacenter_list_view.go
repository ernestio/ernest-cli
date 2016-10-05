package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func printDatacenterList(datacenters []Datacenter) {
	if len(datacenters) == 0 {
		fmt.Println("There are no datacenters created yet.")
		return
	}

	var aws []Datacenter
	var vcloud []Datacenter

	for _, d := range datacenters {
		switch d.Type {
		case "aws", "aws-fake":
			aws = append(aws, d)
		case "vcloud", "vcloud-fake":
			vcloud = append(vcloud, d)
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
}
