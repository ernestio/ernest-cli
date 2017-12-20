/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
	"github.com/urfave/cli"
)

// ListProjects ...
var ListProjects = cli.Command{
	Name:        "list",
	Usage:       h.T("project.list.usage"),
	ArgsUsage:   h.T("project.list.args"),
	Description: h.T("project.list.description"),
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		projects := client.Project().List()
		view.PrintProjectList(projects)

		return nil
	},
}

// InfoProject ...
var InfoProject = cli.Command{
	Name:        "info",
	Usage:       h.T("project.info.usage"),
	ArgsUsage:   h.T("project.info.args"),
	Description: h.T("project.info.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "project.info.args")
		client := esetup(c, AuthUsersValidation)
		p := client.Project().Get(c.Args()[0])
		view.PrintProjectInfo(p)

		return nil
	},
}

// UpdateProjects : Will update the project specific fields
var UpdateProjects = cli.Command{
	Name:        "update",
	Usage:       "Updates an existing project.",
	Description: "Update an existing project on the targeted instance of Ernest.",
	Subcommands: []cli.Command{
		UpdateVCloudProject,
		UpdateAWSProject,
		UpdateAzureProject,
	},
}

// CreateProjects ...
var CreateProjects = cli.Command{
	Name:        "create",
	Usage:       "Create a new project.",
	Description: "Create a new project on the targeted instance of Ernest.",
	Subcommands: []cli.Command{
		CreateVcloudProject,
		CreateAWSProject,
		CreateAzureProject,
	},
}

// CmdProject ...
var CmdProject = cli.Command{
	Name:  "project",
	Usage: "Project related subcommands",
	Subcommands: []cli.Command{
		ListProjects,
		CreateProjects,
		UpdateProjects,
		DeleteProject,
		InfoProject,
	},
}
