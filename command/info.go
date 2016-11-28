/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import "github.com/urfave/cli"

// Info command
// Shows the current ernest target instance information
var Info = cli.Command{
	Name:      "info",
	Aliases:   []string{"i"},
	ArgsUsage: " ",
	Usage:     "Display system-wide information.",
	Description: `Displays ernest instance information.

   Example:
    $ ernest info
	`,
	Action: func(c *cli.Context) error {
		_, cfg := setup(c)
		println("Current target : " + cfg.URL)
		println("Current user : " + cfg.User)
		return nil
	},
}
