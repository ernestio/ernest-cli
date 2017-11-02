/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import "github.com/urfave/cli"

// MapEnvOptions : maps environment options from cli context
func MapEnvOptions(c *cli.Context) map[string]interface{} {
	opts := make(map[string]interface{})

	opts["sync_interval"] = c.String("sync_interval")

	return opts
}
