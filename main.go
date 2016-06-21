/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"os"

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
		Target,
		Info,
		Login,
		CmdUser,
		CmdGroup,
		CmdDatacenter,
		CmdService,
	}
	app.Run(os.Args)
}
