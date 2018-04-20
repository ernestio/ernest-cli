/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"github.com/urfave/cli"
)

// MapEnvOptions : maps environment options from cli context
func MapEnvOptions(c *cli.Context, opts map[string]interface{}) map[string]interface{} {
	if opts == nil {
		opts = make(map[string]interface{})
	}

	submissions := c.String("submissions")

	// default submissions to true
	if !c.IsSet("submissions") && opts["submissions"] == nil {
		opts["submissions"] = true
	}

	if submissions == "enable" {
		opts["submissions"] = true
	}

	if submissions == "disable" {
		opts["submissions"] = false
	}

	return opts
}
