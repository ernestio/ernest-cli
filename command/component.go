/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"

	"github.com/ernestio/ernest-cli/view"
)

// FindComponents ...
var FindComponents = cli.Command{
	Name:      "list",
	Usage:     "List components on your datacenter.",
	ArgsUsage: " ",
	Description: `List all components on your datacenter.

   Example:
    $ ernest component list ebs my_datacenter --service=my_service
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service",
			Value: "",
			Usage: "You can filter by service",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify an existing datacenter name")
			return nil
		}

		if len(c.Args()) == 1 {
			color.Red("You should specify the component type")
			return nil
		}

		datacenter := c.Args()[0]
		component := c.Args()[1]
		service := c.String("service")
		components, err := m.FindComponents(cfg.Token, datacenter, component, service)
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
