package view

import (
	"fmt"
	"os"
	"strings"

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
		table.SetHeader([]string{"Policy Name", "Attached Environments"})
		for _, s := range policies {
			table.Append([]string{s.Name, strings.Join(s.Environments, ",")})
		}
		table.Render()
	}
}
