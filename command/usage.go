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
		cli.StringFlag{
			Name:  "from",
			Usage: "the from date the report will be calculated from. Format YYYY-MM-DD",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "the to date the report will be caluclutated to. Format YYYY-MM-DD",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "the file path to store the report",
		},
	},
	Action: func(c *cli.Context) error {
		var body string
		var err error

		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if body, err = m.GetUsageReport(cfg.Token, c.String("from"), c.String("to")); err != nil {
			h.PrintError(err.Error())
		}

		if c.String("output") != "" {
			if err := ioutil.WriteFile(c.String("output"), []byte(body), 0644); err != nil {
				h.PrintError(err.Error())
			}
			color.Green("A file named " + c.String("output") + " has been exported to the current folder")
		} else {
			fmt.Println(body)
		}

		return nil
	},
}
