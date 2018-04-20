/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"fmt"

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
	n, err := c.cli.Notifications.Get(notification)
	if err != nil {
		h.PrintError(err.Error())
	}

	for _, source := range n.Sources {
		if source == project {
			h.PrintError("project is already added to notification")
		}
	}

	n.Sources = append(n.Sources, project)

	if err := c.cli.Notifications.Update(n); err != nil {
		h.PrintError(err.Error())
	}
}

// RmProject : Removes a project from a notification
func (c *Notification) RmProject(notification, project string) {
	n, err := c.cli.Notifications.Get(notification)
	if err != nil {
		h.PrintError(err.Error())
	}

	for i := len(n.Sources) - 1; i >= 0; i-- {
		if n.Sources[i] == project {
			n.Sources = append(n.Sources[:i], n.Sources[i+1:]...)
		}
	}

	if err := c.cli.Notifications.Update(n); err != nil {
		h.PrintError(err.Error())
	}
}

// AddEnv : Adds an environment to a notification
func (c *Notification) AddEnv(notification, project, env string) {
	n, err := c.cli.Notifications.Get(notification)
	if err != nil {
		h.PrintError(err.Error())
	}

	name := fmt.Sprintf("%s/%s", project, env)

	for _, source := range n.Sources {
		if source == name {
			h.PrintError("environment is already added to notification")
		}
	}

	n.Sources = append(n.Sources, name)

	if err := c.cli.Notifications.Update(n); err != nil {
		h.PrintError(err.Error())
	}
}

// RmEnv : Removes an environment from a notification
func (c *Notification) RmEnv(notification, project, env string) {
	n, err := c.cli.Notifications.Get(notification)
	if err != nil {
		h.PrintError(err.Error())
	}

	name := fmt.Sprintf("%s/%s", project, env)

	for i := len(n.Sources) - 1; i >= 0; i-- {
		if n.Sources[i] == name {
			n.Sources = append(n.Sources[:i], n.Sources[i+1:]...)
		}
	}

	if err := c.cli.Notifications.Update(n); err != nil {
		h.PrintError(err.Error())
	}
}
