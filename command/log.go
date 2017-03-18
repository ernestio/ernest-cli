/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"runtime"

	"github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/nu7hatch/gouuid"
	"github.com/urfave/cli"
)

// CmdLog : Preferences setup
var CmdLog = cli.Command{
	Name:      "log",
	Usage:     "Inline display of ernest logs.",
	ArgsUsage: " ",
	Description: `Display ernest server logs inline

   Example:
    $ ernest log
    $ ernest log --raw
	`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "raw",
			Usage: "Raw output will be displayed instead of pretty-printed",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		uu, _ := uuid.NewV4()
		logger := model.Logger{
			Type: "sse",
			UUID: uu.String(),
		}

		if err := m.SetLogger(cfg.Token, logger); err != nil {
			color.Red(err.Error())
			return nil
		}

		if c.Bool("raw") {
			helper.PrintRawLogs(cfg.URL, "/logs", cfg.Token, logger.UUID)
		} else {
			helper.PrintLogs(cfg.URL, "/logs", cfg.Token, logger.UUID)
		}

		defer func() {
			if err := m.DelLogger(cfg.Token, logger); err != nil {
				color.Red("Ernest wasn't able to reset sse logger")
			}
		}()

		runtime.Goexit()

		return nil
	},
}
