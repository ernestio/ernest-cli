/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// Logout command
// Clear local authentication credentials
var Logout = cli.Command{
	Name:      "logout",
	Usage:     "Clear local authentication credentials.",
	ArgsUsage: " ",
	Description: `Logs out an user from Ernest instance.

   Example:
    $ ernest logout
  `,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if m == nil {
			os.Exit(1)
		}
		cfg.Token = ""
		cfg.User = ""
		err := saveConfig(cfg)
		if err != nil {
			color.Red("Can't write config file")
			os.Exit(1)
		}
		color.Green("Log out succesful.")
		return nil
	},
}
