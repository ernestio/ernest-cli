/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Role : ernest-go-sdk Roles wrapper
type Role struct {
	cli *eclient.Client
}

// Create : Creates a new role
func (c *Role) Create(role *emodels.Role) {
	if err := c.cli.Roles.Create(role); err != nil {
		h.PrintError(err.Error())
	}
}

// Delete : Deletes a role and all its relations
func (c *Role) Delete(role *emodels.Role) {
	if err := c.cli.Roles.Delete(role); err != nil {
		h.PrintError(err.Error())
	}
}
