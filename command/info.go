/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"fmt"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

// Info command
// Shows the current ernest target instance information
var Info = cli.Command{
	Name:        "info",
	Aliases:     []string{"i"},
	Usage:       h.T("info.usage"),
	ArgsUsage:   h.T("info.args"),
	Description: h.T("info.description"),
	Action: func(c *cli.Context) error {
		_, cfg := setup(c)
		fmt.Println("Target:      " + cfg.URL)
		fmt.Println("User:        " + cfg.User)
		fmt.Println("CLI Version: " + c.App.Version)

		return nil
	},
}
