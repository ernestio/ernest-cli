/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"os"

	"github.com/ernestio/ernest-cli/command"
	"github.com/urfave/cli"
)

// Version is the version number
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "ernest"
	app.Version = Version
	app.Usage = "Command line interface for Ernest"
	app.Commands = []cli.Command{
		command.Target,
		command.Info,
		command.Login,
		command.Logout,
		command.CmdUser,
		command.CmdGroup,
		command.CmdDatacenter,
		command.CmdService,
		command.CmdPreferences,
		command.CmdDocs,
		command.CmdSetup,
		command.CmdComponents,
	}
	app.Run(os.Args)
}
