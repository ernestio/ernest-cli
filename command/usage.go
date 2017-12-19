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
		stringFlagND("from", "the from date the report will be calculated from. Format YYYY-MM-DD"),
		stringFlagND("to", "the to date the report will be caluclutated to. Format YYYY-MM-DD"),
		stringFlagND("output", "the file path to store the report"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		body := client.Report().Usage(c.String("from"), c.String("to"))

		if c.String("output") != "" {
			if err := ioutil.WriteFile(c.String("output"), body, 0644); err != nil {
				h.PrintError(err.Error())
			}
			color.Green("A file named " + c.String("output") + " has been exported to the current folder")
		} else {
			fmt.Println(string(body))
		}

		return nil
	},
}
