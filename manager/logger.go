/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Logger : ernest-go-sdk Logger wrapper
type Logger struct {
	cli *eclient.Client
}

// Create : Creates a new logger
func (c *Logger) Create(logger *emodels.Logger) {
	if err := c.cli.Loggers.Create(logger); err != nil {
		h.PrintError(err.Error())
	}
}

// List : lists all available loggers
func (c *Logger) List() []*emodels.Logger {
	loggers, err := c.cli.Loggers.List()
	if err != nil {
		h.PrintError(err.Error())
		return nil
	}
	return loggers
}

// Delete : Deletes logger by name
func (c *Logger) Delete(name string) {
	if err := c.cli.Loggers.Delete(name); err != nil {
		h.PrintError(err.Error())
	}
}
