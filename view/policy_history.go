package view

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintPolicyHistory : Pretty print for policy model
func PrintPolicyHistory(documents []*emodels.PolicyDocument) {
	if len(documents) == 0 {
		fmt.Println("\nThere are no policies created yet")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Revision", "User", "Created"})
		for _, d := range documents {
			table.Append([]string{strconv.Itoa(d.Revision), d.Username, d.CreatedAt})
		}
		table.Render()
	}
}
