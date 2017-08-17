/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"

	"github.com/ernestio/ernest-cli/helper"
)

// NullWriter to disable logging
type NullWriter int

// Write sends to nowhere the log messages
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// MonitorEnv command
// Monitorizes an service and shows the actions being performed on it
var MonitorEnv = cli.Command{
	Name:      "monitor",
	Aliases:   []string{"m"},
	Usage:     "Monitor an environment creation.",
	ArgsUsage: "<project_name> <env_name>",
	Description: `Monitors an environment while it is being built by its name.

   Example:
    $ ernest monitor <my_project> <my_env>
	`,
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
			color.Red("You should specify an existing env name")
			return nil
		}

		project := c.Args()[0]
		env := c.Args()[0]
		service, err := m.EnvStatus(cfg.Token, project, env)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

		if service.Status == "done" {
			color.Yellow("Service has been successfully built")
			color.Yellow("You can check its information running `ernest-cli env info " + project + " / " + env + "`")
			return nil
		}

		return helper.Monitorize(cfg.URL, "/events", cfg.Token, service.ID)
	},
}
