/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdDocs subcommand
import (
	"fmt"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

const docURL = "http://docs.ernest.io/documentation/"

// CmdDocs : Open docs in the default browser
var CmdDocs = cli.Command{
	Name:        "docs",
	Usage:       h.T("docs.usage"),
	ArgsUsage:   h.T("docs.args"),
	Description: h.T("docs.description"),
	Action: func(c *cli.Context) error {
		err := open.Run(docURL)
		if err != nil {
			fmt.Println("Visit ernest.io documentation site : " + docURL)
		}
		return nil
	},
}
