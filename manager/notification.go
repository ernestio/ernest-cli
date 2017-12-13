/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Notification : ernest-go-sdk Notification wrapper
type Notification struct {
	cli *eclient.Client
}

// Get : Gets a notification by name
func (c *Notification) Get(id string) *emodels.Notification {
	notification, err := c.cli.Notifications.Get(id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return notification
}

// Update : Updates a notification
func (c *Notification) Update(notification *emodels.Notification) {
	if err := c.cli.Notifications.Update(notification); err != nil {
		h.PrintError(err.Error())
	}
}

// Create : Creates a new notification
func (c *Notification) Create(notification *emodels.Notification) {
	if err := c.cli.Notifications.Create(notification); err != nil {
		h.PrintError(err.Error())
	}
}

// List : Lists all notifications on the system
func (c *Notification) List() []*emodels.Notification {
	notifications, err := c.cli.Notifications.List()
	if err != nil {
		h.PrintError(err.Error())
	}
	return notifications
}

// Delete : Deletes a notification and all its relations
func (c *Notification) Delete(notification string) {
	if err := c.cli.Notifications.Delete(notification); err != nil {
		h.PrintError(err.Error())
	}
}

// AddProject : Adds a project to a notification
func (c *Notification) AddProject(notification, project string) {
	if err := c.cli.Notifications.AddProject(notification, project); err != nil {
		h.PrintError(err.Error())
	}
}

// RmProject : Removes a project from a notification
func (c *Notification) RmProject(notification, project string) {
	if err := c.cli.Notifications.RemoveProject(notification, project); err != nil {
		h.PrintError(err.Error())
	}
}

// AddEnv : Adds an environment to a notification
func (c *Notification) AddEnv(notification, project, env string) {
	if err := c.cli.Notifications.AddEnv(notification, project, env); err != nil {
		h.PrintError(err.Error())
	}
}

// RmEnv : Removes an environment from a notification
func (c *Notification) RmEnv(notification, project, env string) {
	if err := c.cli.Notifications.RemoveEnv(notification, project, env); err != nil {
		h.PrintError(err.Error())
	}
}
