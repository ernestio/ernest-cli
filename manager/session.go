/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

type session struct {
	cli *eclient.Client
}

// Get : ..
func (c *session) Get() *emodels.Session {
	ses, err := c.cli.Sessions.Get()
	if err != nil {
		h.PrintError("You don’t have permissions to perform this action")
	}
	return ses
}
