/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"fmt"
	"io/ioutil"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CmdUsage : Exports an usage report to a file on the current folder
var CmdUsage = cli.Command{
	Name:        "usage",
	Usage:       h.T("usage.usage"),
	ArgsUsage:   h.T("usage.args"),
	Description: h.T("usage.description"),
	Flags: []cli.Flag{
		tStringFlagND("usage.flags.from"),
		tStringFlagND("usage.flags.to"),
		tStringFlagND("usage.flags.output"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		body := client.Report().Usage(c.String("from"), c.String("to"))

		if c.String("output") != "" {
			if err := ioutil.WriteFile(c.String("output"), body, 0644); err != nil {
				h.PrintError(err.Error())
			}
			color.Green(fmt.Sprintf(h.T("usage.success"), c.String("output")))
		} else {
			fmt.Println(string(body))
		}

		return nil
	},
}
