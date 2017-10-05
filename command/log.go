/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"runtime"

	"github.com/ernestio/ernest-cli/helper"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/nu7hatch/gouuid"
	"github.com/urfave/cli"
)

// CmdLog : Preferences setup
var CmdLog = cli.Command{
	Name:        "log",
	Usage:       h.T("log.usage"),
	ArgsUsage:   h.T("log.args"),
	Description: h.T("log.description"),
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "raw",
			Usage: "Raw output will be displayed instead of pretty-printed",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		uu, _ := uuid.NewV4()
		logger := model.Logger{
			Type: "sse",
			UUID: uu.String(),
		}

		if err := m.SetLogger(cfg.Token, logger); err != nil {
			h.PrintError(err.Error())
		}

		if c.Bool("raw") {
			_ = helper.PrintRawLogs(cfg.URL, "/logs", cfg.Token, logger.UUID)
		} else {
			_ = helper.PrintLogs(cfg.URL, "/logs", cfg.Token, logger.UUID)
		}

		defer func() {
			if err := m.DelLogger(cfg.Token, logger); err != nil {
				h.PrintError("Ernest wasn't able to reset sse logger")
			}
		}()

		runtime.Goexit()

		return nil
	},
}
