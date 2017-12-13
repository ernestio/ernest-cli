/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdUser subcommand
import (
	"github.com/fatih/color"
	"github.com/urfave/cli"

	h "github.com/ernestio/ernest-cli/helper"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

func rolesManager(c *cli.Context, set bool) {
	requiredFlags(c, []string{"role", "user", "project"})
	client := esetup(c, AuthUsersValidation)
	rType := "project"
	rID := c.String("project")
	if c.String("environment") != "" {
		rType = "environment"
		rID = c.String("project") + "/" + c.String("environment")
	}

	role := &emodels.Role{
		ID:       rID,
		User:     c.String("user"),
		Role:     c.String("role"),
		Resource: rType,
	}
	if set {
		client.Role().Create(role)
		verb := "own"
		if c.String("role") == "reader" {
			verb = "read"
		}
		color.Green("User '" + c.String("user") + "' has been authorized to " + verb + " resource " + rID)
	} else {
		client.Role().Delete(role)
		color.Green("User '" + c.String("user") + "' has been unauthorized as " + rID + " " + c.String("role"))
	}

}

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
		rolesManager(c, true)
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
		rolesManager(c, false)
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
