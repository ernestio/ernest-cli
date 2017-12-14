/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"github.com/fatih/color"
	"github.com/urfave/cli"

	h "github.com/ernestio/ernest-cli/helper"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// CreateVcloudProject : Creates a VCloud Project
var CreateVcloudProject = cli.Command{
	Name:        "vcloud",
	Usage:       h.T("vcloud.create.usage"),
	ArgsUsage:   h.T("vcloud.create.args"),
	Description: h.T("vcloud.create.description"),
	Flags: []cli.Flag{
		stringFlag("user", "", "Your VCloud valid user name"),
		stringFlag("password", "", "Your VCloud valid password"),
		stringFlag("org", "", "Your vCloud Organization"),
		stringFlag("vdc", "", "Your vCloud vDC"),
		stringFlag("vcloud-url", "", "VCloud URL"),
		stringFlag("template", "", "Project template"),
		boolFlag("fake", "Fake project"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "vcloud.create.args")
		client := esetup(c, AuthUsersValidation)
		creds := parseTemplateFlags(c, map[string]flagDef{
			"vcloud-url": flagDef{typ: "string", mapto: "vcloud_url", req: true},
			"user":       flagDef{typ: "string", req: true},
			"org":        flagDef{typ: "string", req: true},
			"vdc":        flagDef{typ: "string", req: true},
			"password":   flagDef{typ: "string", req: true},
			"fake":       flagDef{typ: "bool", def: false},
		})
		creds["username"] = creds["user"].(string) + "@" + creds["org"].(string)

		rtype := "vcloud"
		if creds["fake"].(bool) {
			rtype = "vcloud-fake"
		}

		p := &emodels.Project{
			Name:        c.Args()[0],
			Type:        rtype,
			Credentials: creds,
		}
		client.Project().Create(p)
		color.Green("Project '" + p.Name + "' successfully created ")

		return nil
	},
}

// DeleteProject : Project deletion command definition
var DeleteProject = cli.Command{
	Name:        "delete",
	Usage:       h.T("vcloud.delete.usage"),
	ArgsUsage:   h.T("vcloud.delete.args"),
	Description: h.T("vcloud.delete.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "vcloud.delete.args")
		client := esetup(c, AuthUsersValidation)

		name := c.Args()[0]
		client.Project().Delete(name)
		color.Green("Project " + name + " successfully removed")

		return nil
	},
}

// UpdateVCloudProject : Updates the specified VCloud project
var UpdateVCloudProject = cli.Command{
	Name:        "vcloud",
	Usage:       h.T("vcloud.update.usage"),
	ArgsUsage:   h.T("vcloud.update.args"),
	Description: h.T("vcloud.update.description"),
	Flags: []cli.Flag{
		stringFlag("user", "", "Your VCloud valid user name"),
		stringFlag("password", "", "Your VCloud valid password"),
		stringFlag("org", "", "Your vCloud Organization"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "vcloud.update.args")
		client := esetup(c, AuthUsersValidation)

		creds := parseTemplateFlags(c, map[string]flagDef{
			"user":     flagDef{typ: "string"},
			"org":      flagDef{typ: "string"},
			"password": flagDef{typ: "string"},
		})
		creds["user"] = creds["user"].(string) + "@" + creds["org"].(string)

		n := client.Project().Get(c.Args()[0])
		n.Credentials["user"] = creds["user"].(string)
		n.Credentials["passwrord"] = creds["password"].(string)
		client.Project().Update(n)
		color.Green("Project " + n.Name + " successfully updated")

		return nil
	},
}
