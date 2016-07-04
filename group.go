/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdUser subcommand
import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListGroups ...
var ListGroups = cli.Command{
	Name:      "list",
	Usage:     "List available groups.",
	ArgsUsage: " ",
	Description: `List available groups.

   Example:
    $ ernest group list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		groups, err := m.ListGroups(cfg.Token)
		if err != nil {
			color.Red(err.Error())
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "NAME")
		for _, group := range groups {
			str := fmt.Sprintf("%s", group.Name)
			fmt.Fprintln(w, str)
		}
		w.Flush()
		return nil
	},
}

// CmdGroup ...
var CmdGroup = cli.Command{
	Name:  "group",
	Usage: "Group related subcommands",
	Subcommands: []cli.Command{
		ListGroups,
	},
}
