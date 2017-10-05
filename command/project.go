/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"errors"
	"io/ioutil"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

// ListProjects ...
var ListProjects = cli.Command{
	Name:        "list",
	Usage:       h.T("project.list.usage"),
	ArgsUsage:   h.T("project.list.args"),
	Description: h.T("project.list.description"),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}
		projects, err := m.ListProjects(cfg.Token)
		if err != nil {
			h.PrintError(err.Error())
		}

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
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}
		if len(c.Args()) == 0 {
			h.PrintError("You should specify the project name")
		}
		project := c.Args()[0]
		p, err := m.InfoProject(cfg.Token, project)
		if err != nil {
			h.PrintError(err.Error())
		}

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

func getProjectTemplate(template string, t *model.ProjectTemplate) (err error) {
	payload, err := ioutil.ReadFile(template)
	if err != nil {
		return errors.New("Template file '" + template + "' not found")
	}
	if yaml.Unmarshal(payload, &t) != nil {
		return errors.New("Template file '" + template + "' is not valid yaml file")
	}
	return err
}
