/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Environment : ernest-go-sdk Environment wrapper
type Environment struct {
	cli *eclient.Client
}

// Create : ...
func (c *Environment) Create(project string, env *emodels.Environment) {
	if err := c.cli.Environments.Create(project, env); err != nil {
		h.PrintError(err.Error())
	}
}

// Delete : Deletes a env and all its relations
func (c *Environment) Delete(project, env string) *emodels.Build {
	build, err := c.cli.Environments.Delete(project, env)
	if err != nil {
		h.PrintError(err.Error())
	}
	return build
}

// ForceDeletion : Deletes a env and all its relations
func (c *Environment) ForceDeletion(project, env string) *emodels.Build {
	build, err := c.cli.Environments.ForceDeletion(project, env)
	if err != nil {
		h.PrintError(err.Error())
	}
	return build
}

// Get : Gets a env by name
func (c *Environment) Get(project, id string) *emodels.Environment {
	env, err := c.cli.Environments.Get(project, id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return env
}

// Sync : Syncs a env by name
func (c *Environment) Sync(project, id string) *emodels.Action {
	act, err := c.cli.Environments.Sync(project, id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return act
}

// Resolve : Resolves a env by name
func (c *Environment) Resolve(project, id, resolution string) *emodels.Action {
	act, err := c.cli.Environments.Resolve(project, id, resolution)
	if err != nil {
		h.PrintError(err.Error())
	}
	return act
}

// Reset : Resets a env by name
func (c *Environment) Reset(project, id string) *emodels.Action {
	act, err := c.cli.Environments.Reset(project, id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return act
}

// Update : Updates a notification
func (c *Environment) Update(env *emodels.Environment) {
	if err := c.cli.Environments.Update(env); err != nil {
		h.PrintError(err.Error())
	}
}

// ListAll : Lists all envs on the system
func (c *Environment) ListAll() []*emodels.Environment {
	envs, err := c.cli.Environments.ListAll()
	if err != nil {
		h.PrintError(err.Error())
	}
	return envs
}

// Import : creates an import build for an environment
func (c *Environment) Import(project, env string, filters []string) *emodels.Action {
	action, err := c.cli.Environments.Import(project, env, filters)
	if err != nil {
		h.PrintError(err.Error())
	}
	return action
}
