/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"net/url"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// Target command
// Configures the ernest target instance
var Target = cli.Command{
	Name:        "target",
	Aliases:     []string{"t"},
	Usage:       h.T("target.usage"),
	ArgsUsage:   h.T("target.args"),
	Description: h.T("target.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "target.args")
		cfg := esetup(c, NoValidation).Config()

		cfg.URL = c.Args()[0]
		persistTarget(cfg)

		color.Green("Target set")
		return nil
	},
}

func persistTarget(cfg *model.Config) {
	u, _ := url.Parse(cfg.URL)
	if u.Scheme == "http" {
		color.Yellow("Warning! Your are using an insecure target for Ernest")
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		color.Red("You should specify a valid url for the target")
		return
	}
	err := model.SaveConfig(cfg)
	if err != nil {
		color.Red(err.Error())
		h.EvaluateErrorMsg(err, "Couldn't write config file ~/.ernest check permissions")
	}
}
