/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdUser subcommand
import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CmdRolesSet :
var CmdRolesSet = cli.Command{
	Name:      "set",
	Usage:     "ernest role set -u john -r owner -p project",
	ArgsUsage: "",
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
	Description: `Set permissions for a user on a specific resource

   Example:
    $ ernest roles set -u john -r owner -p my_project 
    $ ernest roles set -u john -r reader -p my_project -e my_environment
	`,
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

		color.Green("User " + u + " authorized as " + r + " to " + resource)
		return nil
	},
}

// CmdRolesUnset :
var CmdRolesUnset = cli.Command{
	Name:      "unset",
	Usage:     "ernest role unset -u john -r owner -p my_project",
	ArgsUsage: "",
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
	Description: `Set permissions for a user on a specific resource

   Example:
    $ ernest roles set -u john -r owner -p my_project 
    $ ernest roles set -u john -r reader -p my_project -e my_environment
	`,
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

		color.Green("User " + u + " unauthorized as " + r + " to " + resource)
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
