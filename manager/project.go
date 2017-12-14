/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Project : ernest-go-sdk Project wrapper
type Project struct {
	cli *eclient.Client
}

// Create : ...
func (c *Project) Create(project *emodels.Project) {
	if err := c.cli.Projects.Create(project); err != nil {
		h.PrintError(err.Error())
	}
}

// Delete : Deletes a project and all its relations
func (c *Project) Delete(project string) {
	if err := c.cli.Projects.Delete(project); err != nil {
		h.PrintError(err.Error())
	}
}

// Get : Gets a project by name
func (c *Project) Get(id string) *emodels.Project {
	project, err := c.cli.Projects.Get(id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return project
}

// Update : Updates a notification
func (c *Project) Update(project *emodels.Project) {
	if err := c.cli.Projects.Update(project); err != nil {
		h.PrintError(err.Error())
	}
}

// List : Lists all projects on the system
func (c *Project) List() []*emodels.Project {
	projects, err := c.cli.Projects.List()
	if err != nil {
		h.PrintError(err.Error())
	}
	return projects
}
