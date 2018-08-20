package view

import (
	"fmt"
	"os"
	"strings"

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
	table.SetHeader([]string{"ID", "Action", "Interval", "Instance", "Resolution"})
	for k, d := range list {
		var instances []string

		opts := d.(map[string]interface{})
		t := opts["type"].(string)

		if len(opts["instances"].([]interface{})) != 0 {
			for _, v := range opts["instances"].([]interface{}) {
				instances = append(instances, v.(string))
			}
		}

		resolution, _ := opts["resolution"].(string)

		i := opts["interval"].(string)
		table.Append([]string{k, t, i, strings.Join(instances, ","), resolution})
	}
	table.Render()
}
