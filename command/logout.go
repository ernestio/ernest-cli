/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"os"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// Logout command
// Clear local authentication credentials
var Logout = cli.Command{
	Name:        "logout",
	Usage:       h.T("logout.usage"),
	ArgsUsage:   h.T("logout.args"),
	Description: h.T("logout.description"),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're already logged out")
		}
		if m == nil {
			os.Exit(1)
		}
		cfg.Token = ""
		cfg.User = ""
		err := model.SaveConfig(cfg)
		if err != nil {
			h.PrintError("Can't write config file")
		}
		color.Green("Bye.")
		return nil
	},
}
