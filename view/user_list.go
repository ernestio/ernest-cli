package view

import (
	"os"
	"strconv"
	"text/tabwriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
	"github.com/olekukonko/tablewriter"
)

// PrintUserList : ...
func PrintUserList(users []*emodels.User) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Type", "Admin", "Disabled"})
	for _, u := range users {
		id := strconv.Itoa(u.ID)
		admin := "no"
		if u.Admin {
			admin = "yes"
		}
		disabled := "no"
		if u.Disabled {
			disabled = "yes"
		}
		table.Append([]string{id, u.Username, u.Type, admin, disabled})
	}
	table.Render()

}
