/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
)

// FindComponents ...
var FindComponents = cli.Command{
	Name:        "list",
	Usage:       h.T("components.find.usage"),
	Description: h.T("components.find.description"),
	ArgsUsage:   h.T("components.find.args"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "environment",
			Value: "",
			Usage: "You can filter by environment",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify an existing project name")
			return nil
		}

		if len(c.Args()) == 1 {
			color.Red("You should specify the component type")
			return nil
		}

		project := c.Args()[0]
		component := c.Args()[1]
		service := c.String("environment")
		components, err := m.FindComponents(cfg.Token, project, component, service)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		view.PrintComponentsList(components)

		return nil
	},
}

// CmdComponents ...
var CmdComponents = cli.Command{
	Name:  "component",
	Usage: "Components related subcommands",
	Subcommands: []cli.Command{
		FindComponents,
	},
}
