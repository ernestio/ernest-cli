/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
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
		client := esetup(c, AuthUsersValidation)

		cfg := client.Config()
		cfg.Token = ""
		cfg.User = ""
		if err := model.SaveConfig(cfg); err != nil {
			h.PrintError(h.T("logout.errors.write"))
		}
		color.Green("faaaaail")
		return nil
	},
}
