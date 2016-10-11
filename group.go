/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdUser subcommand
import (
	"errors"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// DeleteGroup ...
var DeleteGroup = cli.Command{
	Name:      "delete",
	Usage:     "Deletes a group.",
	ArgsUsage: " ",
	Description: `Deletes a group by name

	  Example:
		  $ ernest group delete <name>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the group name"
			color.Red(msg)
			return errors.New(msg)
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		if session.IsAdmin == false {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		name := c.Args()[0]
		err = m.DeleteGroup(cfg.Token, c.Args()[0])
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Group '" + name + "' successfully deleted")
		return nil
	},
}

// RemoveDatacenter ...
var RemoveDatacenter = cli.Command{
	Name:      "remove-datacenter",
	Usage:     "Removes a datacenter from a group.",
	ArgsUsage: " ",
	Description: `Removes an datacenter from a group.

		Example:
		  $ ernest group remove-datacenter <datacenter name> <group name>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the datacenter name and group name"
			color.Red(msg)
			return errors.New(msg)
		}
		if len(c.Args()) < 2 {
			msg := "You should specify the group name"
			color.Red(msg)
			return errors.New(msg)
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		if session.IsAdmin == false {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		datacenter := c.Args()[0]
		group := c.Args()[1]
		err = m.GroupRemoveDatacenter(cfg.Token, datacenter, group)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Datacenter '" + datacenter + "' is not assigned anymore to group '" + group + "'")
		return nil
	},
}

// AddDatacenter ...
var AddDatacenter = cli.Command{
	Name:      "add-datacenter",
	Usage:     "Adds a datacenter to a group.",
	ArgsUsage: " ",
	Description: `Adds a datacenter to a group.

	  Example:
		  $ ernest group add-datacenter <datacenter-name> <group-name>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the datacenter name and group name"
			color.Red(msg)
			return errors.New(msg)
		}
		if len(c.Args()) < 2 {
			msg := "You should specify the group name"
			color.Red(msg)
			return errors.New(msg)
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		if session.IsAdmin == false {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		datacenter := c.Args()[0]
		group := c.Args()[1]
		err = m.GroupAddDatacenter(cfg.Token, datacenter, group)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Datacenter '" + datacenter + "' is now assigned to group '" + group + "'")
		return nil
	},
}

// RemoveUser ...
var RemoveUser = cli.Command{
	Name:      "remove-user",
	Usage:     "Removes an user from a group.",
	ArgsUsage: " ",
	Description: `Removes an user from a group.

		Example:
		  $ ernest group remove-user <username> <group name>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the username and group name"
			color.Red(msg)
			return errors.New(msg)
		}
		if len(c.Args()) < 2 {
			msg := "You should specify the group name"
			color.Red(msg)
			return errors.New(msg)
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		if session.IsAdmin == false {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		user := c.Args()[0]
		group := c.Args()[1]
		err = m.GroupRemoveUser(cfg.Token, user, group)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("User '" + user + "' is not assigned anymore to group '" + group + "'")
		return nil
	},
}

// AddUser : Adds a user to a group
var AddUser = cli.Command{
	Name:      "add-user",
	Usage:     "Adds a user to a group.",
	ArgsUsage: " ",
	Description: `Adds a user to a group.

	  Example:
		  $ ernest group add-user <user-name> <group-name>
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the username and group name"
			color.Red(msg)
			return errors.New(msg)
		}
		if len(c.Args()) < 2 {
			msg := "You should specify the group name"
			color.Red(msg)
			return errors.New(msg)
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		if session.IsAdmin == false {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		user := c.Args()[0]
		group := c.Args()[1]
		err = m.GroupAddUser(cfg.Token, user, group)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("User '" + user + "' is now assigned to group '" + group + "'")
		return nil
	},
}

// CreateGroup : Creates a group
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
			color.Red("You should specify the group name")
			return nil
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		if session.IsAdmin == false {
			color.Red("You don't have permissions to perform this action")
			return nil
		}

		name := c.Args()[0]
		err = m.CreateGroup(cfg.Token, c.Args()[0])
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Group '" + name + "' successfully created, you can add users with 'ernest group add-user username " + name + "'")
		return nil
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
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		groups, err := m.ListGroups(cfg.Token)
		if err != nil {
			color.Red("We didn't found any accessible group")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name"})
		for _, g := range groups {
			id := strconv.Itoa(g.ID)
			table.Append([]string{id, g.Name})
		}
		table.Render()

		return nil
	},
}

// CmdGroup ...
var CmdGroup = cli.Command{
	Name:  "group",
	Usage: "Group related subcommands",
	Subcommands: []cli.Command{
		DeleteGroup,
		ListGroups,
		CreateGroup,
		AddUser,
		RemoveUser,
		AddDatacenter,
		RemoveDatacenter,
	},
}
