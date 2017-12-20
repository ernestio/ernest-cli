/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// User : ernest-go-sdk User wrapper
type User struct {
	cli *eclient.Client
}

// Get : ...
func (c *User) Get(username string) *emodels.User {
	user, err := c.cli.Users.Get(username)
	if err != nil {
		h.PrintError(err.Error())
	}
	return user
}

// Update : ...
func (c *User) Update(user *emodels.User) {
	if err := c.cli.Users.Update(user); err != nil {
		h.PrintError(err.Error())
	}
}

// Create : ...
func (c *User) Create(user *emodels.User) {
	if err := c.cli.Users.Create(user); err != nil {
		h.PrintError(err.Error())
	}
}

// List : ...
func (c *User) List() []*emodels.User {
	users, err := c.cli.Users.List()
	if err != nil {
		h.PrintError(err.Error())
	}
	return users
}

// Promote : ...
func (c *User) Promote(user *emodels.User) {
	user.Admin = true
	if err := c.cli.Users.Update(user); err != nil {
		str1 := "It was not possible to set this user as admin: "
		str2 := "Please fix any errors and try again with 'user admin add ...' command"
		h.PrintError(str1 + err.Error() + "\n" + str2)
	}
}

// ToggleMFA : ...
func (c *User) ToggleMFA(user *emodels.User, toggle bool) (res string) {
	user.MFA = &toggle
	c.Update(user)

	if toggle {
		res = user.MFASecret
	}
	return
}
