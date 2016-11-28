/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"github.com/ernestio/ernest-cli/manager"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// setup ...
func setup(c *cli.Context) (*manager.Manager, *model.Config) {
	config := model.GetConfig()
	if config == nil {
		config = &model.Config{}
		if c.Command.Name != "target" {
			color.Red("Environment not configured, please use target command")
		}
	}
	m := manager.Manager{URL: config.URL}
	return &m, config
}
