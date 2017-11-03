/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import "github.com/urfave/cli"

// MapEnvOptions : maps environment options from cli context
func MapEnvOptions(c *cli.Context, opts map[string]interface{}) map[string]interface{} {
	if opts == nil {
		opts = make(map[string]interface{})
	}

	syncInterval := c.String("sync_interval")

	if syncInterval != "" {
		opts["sync_interval"] = syncInterval
	}

	submissions := c.String("submissions")

	if submissions == "enable" {
		opts["submissions"] = true
	}

	if submissions == "disable" {
		opts["submissions"] = false
	}

	return opts
}
