package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func printServiceInfo(service *Service) {
	fmt.Println("Name : " + service.Name)
	if service.Vpc != "" {
		fmt.Println("VPC : " + service.Vpc)
	}
	fmt.Println("Status : " + service.Status)
	if service.Status == "errored" {
		if service.LastError == "" {
			fmt.Println("Last known error : unknown")
		} else {
			fmt.Println("Last known error : " + service.LastError)
		}
	}

	if len(service.ELBs) == 0 {
		fmt.Println("\nELBs (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nELBs:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "DNS Name"})
		for _, v := range service.ELBs {
			table.Append([]string{v.Name, v.DNSName})
		}
		table.Render()
	}

	if len(service.Networks) == 0 {
		fmt.Println("\nNetworks (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nNetworks:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ID", "Availability Zone"})
		for _, v := range service.Networks {
			table.Append([]string{v.Name, v.Subnet, v.AvailabilityZone})
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

	if len(service.Nats) == 0 {
		fmt.Println("\nNAT gateways (empty)")
		fmt.Println("")
	} else {
		fmt.Println("\nNAT gateways:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Group ID"})
		for _, v := range service.Nats {
			table.Append([]string{v.Name, v.NatGatewayAWSID})
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
			table.Append([]string{v.Name, v.SecurityGroupAWSID})
		}
		table.Render()
	}
}
