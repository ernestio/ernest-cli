/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdUser subcommand
import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// AddUser ...
var AddUser = cli.Command{
	Name:      "add-user",
	Usage:     "Adds a user to a group.",
	ArgsUsage: " ",
	Description: `Adds a user to a group.

	  Example:
		  $ ernest group add-user <user-id> <group-id>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 2 {
			msg := "You should specify an user id and a group id."
			color.Red(msg)
			return errors.New(msg)
		}
		m, cfg := setup(c)
		userid := c.Args()[0]
		groupid := c.Args()[1]
		err := m.GroupAddUser(cfg.Token, userid, groupid)
		return err
	},
}

// CreateGroup ...
var CreateGroup = cli.Command{
	Name:      "create",
	Usage:     "Create a group.",
	ArgsUsage: " ",
	Description: `Create a group.

   Example:
    $ ernest group create <name>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the group name"
			color.Red(msg)
			return errors.New(msg)
		}
		m, cfg := setup(c)
		name := c.Args()[0]
		err := m.CreateGroup(cfg.Token, name)
		return err
	},
}

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
		fmt.Fprintln(w, "NAME\tID")
		for _, group := range groups {
			str := fmt.Sprintf("%d\t%s", group.ID, group.Name)
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
		CreateGroup,
		AddUser,
	},
}
