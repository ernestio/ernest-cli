/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"

	"github.com/ernestio/ernest-cli/model"
	"github.com/olekukonko/tablewriter"
)

// PrintEnvInfo : Pretty print for build info
func PrintEnvInfo(build *model.Build) {
	fmt.Println("Name : " + build.Name)
	fmt.Println("Status : " + build.Status)
	fmt.Println("Project : " + build.ProjectName)
	fmt.Println("Provider : ")
	fmt.Println("  Type : " + build.Provider)
	fmt.Println("Members:")
	for _, m := range build.Roles {
		fmt.Println("  " + m)
	}
	fmt.Println("Date : " + build.CreatedAt)

	if len(build.VPCs) > 0 {
		fmt.Println("\nVPCs:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Subnet"})
		for _, v := range build.VPCs {
			table.Append([]string{v.Name, v.ID, v.Subnet})
		}
		table.Render()
	}

	if len(build.ELBs) > 0 {
		fmt.Println("\nELBs:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "DNS Name"})
		for _, v := range build.ELBs {
			table.Append([]string{v.Name, v.DNSName})
		}
		table.Render()
	}

	if len(build.Networks) > 0 {
		fmt.Println("\nNetworks:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Availability Zone"})
		for _, v := range build.Networks {
			table.Append([]string{v.Name, v.Subnet, v.AvailabilityZone})
		}
		table.Render()
	}

	if len(build.Instances) > 0 {
		fmt.Println("\nInstances:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Public IP", "Private IP"})
		for _, v := range build.Instances {
			table.Append([]string{v.Name, v.InstanceAWSID, v.PublicIP, v.IP})
		}
		table.Render()
	}

	if len(build.Nats) > 0 {
		fmt.Println("\nNAT gateways:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "IP"})
		for _, v := range build.Nats {
			table.Append([]string{v.Name, v.NatGatewayAWSID, v.IP})
		}
		table.Render()
	}

	if len(build.SecurityGroups) > 0 {
		fmt.Println("\nSecurity groups:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Group ID"})
		for _, v := range build.SecurityGroups {
			table.Append([]string{v.Name, v.SecurityGroupAWSID})
		}
		table.Render()
	}

	if len(build.RDSClusters) > 0 {
		fmt.Println("\nRDS Clusters:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Endpoint"})
		for _, v := range build.RDSClusters {
			table.Append([]string{v.Name, v.Endpoint})
		}
		table.Render()
	}

	if len(build.RDSInstances) > 0 {
		fmt.Println("\nRDS Instances:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Endpoint"})
		for _, v := range build.RDSInstances {
			table.Append([]string{v.Name, v.Endpoint})
		}
		table.Render()
	}

	if len(build.EBSVolumes) > 0 {
		fmt.Println("\nEBS Volumes:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Volume ID"})
		for _, v := range build.EBSVolumes {
			table.Append([]string{v.Name, v.VolumeAWSID})
		}
		table.Render()
	}

	if len(build.LoadBalancers) > 0 {
		fmt.Println("\nLoad Balancers:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "IP"})
		for _, v := range build.LoadBalancers {
			table.Append([]string{v.Name, v.PublicIP})
		}
		table.Render()
	}

	if len(build.VirtualMachines) > 0 {
		fmt.Println("\nVirtual Machines:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Public IP", "Private IP"})
		for _, v := range build.VirtualMachines {
			table.Append([]string{v.Name, v.PublicIP, v.PrivateIP})
		}
		table.Render()
	}

	if len(build.SQLDatabases) > 0 {
		fmt.Println("\nSQL Databases:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Server Name"})
		for _, v := range build.SQLDatabases {
			table.Append([]string{v.Name, v.ServerName})
		}
		table.Render()
	}

}
