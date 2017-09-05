/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdUser subcommand
import (
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CmdRolesSet :
var CmdRolesSet = cli.Command{
	Name:        "set",
	Usage:       h.T("roles.set.usage"),
	ArgsUsage:   h.T("roles.set.args"),
	Description: h.T("roles.set.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: "",
			Usage: "User to be authorized over the given resource",
		},
		cli.StringFlag{
			Name:  "project, p",
			Value: "",
			Usage: "Project to authorize",
		},
		cli.StringFlag{
			Name:  "role, r",
			Value: "",
			Usage: "Role type [owner, reader]",
		},
		cli.StringFlag{
			Name:  "environment, e",
			Value: "",
			Usage: "Environment to authorize",
		},
	},
	Action: func(c *cli.Context) error {
		r := c.String("role")
		u := c.String("user")
		p := c.String("project")
		e := c.String("environment")
		if r == "" {
			color.Red("Please provide a role with --role flag")
			return nil
		}
		if u == "" {
			color.Red("Please provide a user with --user flag")
			return nil
		}
		if p == "" {
			color.Red("Please provide a project with --project flag")
			return nil
		}

		m, cfg := setup(c)
		body, err := m.SetRole(cfg.Token, u, p, e, r)
		if err != nil {
			color.Red(body)
			return nil
		}
		resource := p
		if e != "" {
			resource = p + " / " + e
		}
		verb := "own"
		if r == "reader" {
			verb = "read"
		}
		color.Green("User '" + u + "' has been authorized to " + verb + " resource " + resource)

		return nil
	},
}

// CmdRolesUnset :
var CmdRolesUnset = cli.Command{
	Name:        "unset",
	Usage:       h.T("roles.unset.usage"),
	ArgsUsage:   h.T("roles.unset.args"),
	Description: h.T("roles.unset.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: "",
			Usage: "User to be authorized over the given resource",
		},
		cli.StringFlag{
			Name:  "project, p",
			Value: "",
			Usage: "Project to authorize",
		},
		cli.StringFlag{
			Name:  "role, r",
			Value: "",
			Usage: "Role type [owner, reader]",
		},
		cli.StringFlag{
			Name:  "environment, e",
			Value: "",
			Usage: "Environment to authorize",
		},
	},
	Action: func(c *cli.Context) error {
		r := c.String("role")
		u := c.String("user")
		p := c.String("project")
		e := c.String("environment")
		if r == "" {
			color.Red("Please provide a role with --role flag")
			return nil
		}
		if u == "" {
			color.Red("Please provide a user with --user flag")
			return nil
		}
		if p == "" {
			color.Red("Please provide a project with --project flag")
			return nil
		}

		m, cfg := setup(c)
		body, err := m.UnsetRole(cfg.Token, u, p, e, r)
		if err != nil {
			color.Red(body)
			return nil
		}

		resource := p
		if e != "" {
			resource = p + " / " + e
		}
		color.Green("User '" + u + "' has been unauthorized as " + resource + " " + r)

		return nil
	},
}

// CmdRoles ...
var CmdRoles = cli.Command{
	Name:  "role",
	Usage: "Roles to manage resources authorization",
	Subcommands: []cli.Command{
		CmdRolesSet,
		CmdRolesUnset,
	},
}
