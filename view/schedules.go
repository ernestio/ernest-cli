package view

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// PrintScheduleList : ..
func PrintScheduleList(list map[string]interface{}) {
	fmt.Println("")
	if len(list) == 0 {
		fmt.Println("There are no schedules created for this environment")
		fmt.Println("please use 'ernest env schedule add' to create a new one")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Action", "Instance Type", "Interval"})
	for k, d := range list {
		opts := d.(map[string]interface{})
		t := opts["action"].(string)
		it := opts["instance_type"].(string)
		i := opts["interval"].(string)
		table.Append([]string{k, t, it, i})
	}
	table.Render()
}
