/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
)

// Report : ernest-go-sdk Report wrapper
type Report struct {
	cli *eclient.Client
}

// Usage : Gets an usage report
func (c *Report) Usage(from, to string) []byte {
	usage, err := c.cli.Reports.Usage(from, to)
	if err != nil {
		h.PrintError(err.Error())
	}
	return usage
}
