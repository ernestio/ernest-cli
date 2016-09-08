package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printServiceInfo(service *Service) {
	fmt.Println("Name : " + service.Name)
	fmt.Println("Datacenter : " + service.Name)
	fmt.Println("Service IP : " + service.Endpoint)

	if len(service.Networks) == 0 {
		fmt.Println("\nNetworks (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nNetworks:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID"})
		for _, v := range service.Networks {
			table.Append([]string{v.Name, v.Subnet})
		}
		table.Render()
	}

	if len(service.Instances) == 0 {
		fmt.Println("\nInstances (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nInstances:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Public IP", "Private IP"})
		for _, v := range service.Instances {
			table.Append([]string{v.Name, v.InstanceAWSID, v.PublicIP, v.IP})
		}
		table.Render()
	}

	if len(service.SecurityGroups) == 0 {
		fmt.Println("\nSecurity groups (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nSecurity groups:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Group ID"})
		for _, v := range service.SecurityGroups {
			table.Append([]string{v.Name, v.NatGatewayAWSID})
		}
		table.Render()
	}
}
