/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"strconv"

	"github.com/r3labs/diff"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Build : ernest-go-sdk Build wrapper
type Build struct {
	cli     *eclient.Client
	Verbose bool
}

// Create : Creates a new build
func (c *Build) Create(definition []byte) *emodels.Build {
	build, err := c.cli.Builds.Create(definition)
	if err != nil {
		merr, ok := err.(*emodels.Error)
		if ok {
			view.PrintValidation(merr.Validation)
		}
		h.PrintError(err.Error())
	}
	if c.Verbose {
		view.PrintValidation(build.Validation)
	}
	return build
}

// Dry : Simulates the creation of a new build
func (c *Build) Dry(definition []byte) *[]string {
	build, err := c.cli.Builds.Dry(definition)
	if err != nil {
		h.PrintError(err.Error())
	}
	return build
}

// Get : Gets a build by name
func (c *Build) Get(project, env, id string) *emodels.Build {
	build, err := c.cli.Builds.Get(project, env, id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return build
}

// List : Lists all builds on the system
func (c *Build) List(project, env string) []*emodels.Build {
	builds, err := c.cli.Builds.List(project, env)
	if err != nil {
		h.PrintError(err.Error())
	}
	return builds
}

// Diff : Diff two builds by id
func (c *Build) Diff(project, env, from, to string) *diff.Changelog {
	changelog, err := c.cli.Builds.Diff(project, env, from, to)
	if err != nil {
		h.PrintError(err.Error())
	}
	return changelog
}

// Changelog : get a changelog for a build (if it has been generated)
func (c *Build) Changelog(project, env, id string) *diff.Changelog {
	changelog, err := c.cli.Builds.Changelog(project, env, id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return changelog
}

// Stream : Streams build progress
func (c *Build) Stream(id string) chan []byte {
	ch, err := c.cli.Builds.Stream(id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return ch
}

// BuildByPosition : Streams build progress
func (c *Build) BuildByPosition(project, env, pos string) *emodels.Build {
	builds := c.List(project, env)
	if len(builds) == 0 {
		h.PrintError("No builds were found for the specified parameters")
	}

	num := 0

	if pos != "" {
		num, _ = strconv.Atoi(pos)
		if num < 1 || num > len(builds) {
			h.PrintError("Invalid build ID")
		}
		num = len(builds) - num
	}

	return builds[num]
}

// Definition : Gets a build definitin by name
func (c *Build) Definition(project, env, id string) string {
	build, err := c.cli.Builds.Definition(project, env, id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return build
}
