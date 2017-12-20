/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"

	h "github.com/ernestio/ernest-cli/helper"
)

// NullWriter to disable logging
type NullWriter int

// Write sends to nowhere the log messages
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// MonitorEnv command
// Monitorizes an environment and shows the actions being performed on it
var MonitorEnv = cli.Command{
	Name:        "monitor",
	Aliases:     []string{"m"},
	Usage:       h.T("monitor.usage"),
	ArgsUsage:   h.T("monitor.args"),
	Description: h.T("monitor.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "monitor.args")
		client := esetup(c, AuthUsersValidation)

		build := client.Build().BuildByPosition(c.Args()[0], c.Args()[1], "")

		if build.Status == "done" {
			color.Yellow(h.T("monitor.success_1"))
			color.Yellow(fmt.Sprintf(h.T("monitor.success_2"), c.Args()[0], c.Args()[1]))
			return nil
		}
		return h.Monitorize(client.Build().Stream(build.ID))
	},
}
