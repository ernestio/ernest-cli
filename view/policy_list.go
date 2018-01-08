package view

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintPolicyList : Pretty print for policy model
func PrintPolicyList(policies []*emodels.Policy) {
	if len(policies) == 0 {
		fmt.Println("\nThere are no policies created yet")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name"})
		for _, s := range policies {
			table.Append([]string{s.Name})
		}
		table.Render()
	}
}
