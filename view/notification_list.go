package view

import (
	"fmt"
	"os"

	"github.com/ernestio/ernest-cli/model"
	"github.com/olekukonko/tablewriter"
)

// PrintNotificationList : Pretty print for notification model
func PrintNotificationList(notifications []model.Notification) {
	if len(notifications) == 0 {
		fmt.Println("\nThere are no notifications created yet")
		fmt.Println("")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Type", "Config", "Members"})
		for _, s := range notifications {
			table.Append([]string{s.Name, s.Type, s.Config, s.Members})
		}
		table.Render()
	}
}
