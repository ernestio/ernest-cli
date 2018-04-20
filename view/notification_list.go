package view

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintNotificationList : Pretty print for notification model
func PrintNotificationList(notifications []*emodels.Notification) {
	if len(notifications) == 0 {
		fmt.Println("\nThere are no notifications created yet")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Type", "Config", "Members"})
		for _, s := range notifications {
			table.Append([]string{s.Name, s.Type, s.Config, strings.Join(s.Sources, ", ")})
		}
		table.Render()
	}
}
