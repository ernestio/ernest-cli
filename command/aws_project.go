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

// CreateAWSProject : Creates an AWS project
var CreateAWSProject = cli.Command{
	Name:        "aws",
	Usage:       h.T("aws.create.usage"),
	Description: h.T("aws.create.description"),
	ArgsUsage:   h.T("aws.create.args"),
	Flags: []cli.Flag{
		stringFlag("region, r", "", "Project region"),
		stringFlag("access_key_id, k", "", "AWS access key id"),
		stringFlag("secret_access_key, s", "", "AWS Secret access key"),
		stringFlag("template, t", "", "Project template"),
		boolFlag("fake, f", "Fake project"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "aws.create.args")
		client := esetup(c, AuthUsersValidation)
		creds := parseTemplateFlags(c, map[string]flagDef{
			"secret_access_key": flagDef{typ: "string", mapto: "aws_secret_access_key", req: true},
			"access_key_id":     flagDef{typ: "string", mapto: "aws_access_key_id", req: true},
			"region":            flagDef{typ: "string", req: true},
			"fake":              flagDef{typ: "bool", def: false},
		})
		creds["username"] = c.Args()[0]

		rtype := "aws"
		if creds["fake"].(bool) {
			rtype = "aws-fake"
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

// UpdateAWSProject : Updates the specified VCloud project
var UpdateAWSProject = cli.Command{
	Name:        "aws",
	Usage:       h.T("aws.create.usage"),
	ArgsUsage:   h.T("aws.create.args"),
	Description: h.T("aws.create.description"),
	Flags: []cli.Flag{
		stringFlag("access_key_id", "", "Your AWS access key id"),
		stringFlag("secret_access_key", "", "Your AWS secret access key"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "aws.create.args")
		client := esetup(c, AuthUsersValidation)

		creds := parseTemplateFlags(c, map[string]flagDef{
			"access_key_id":     flagDef{typ: "string", mapto: "aws_access_key_id"},
			"secret_access_key": flagDef{typ: "string", mapto: "aws_secret_access_key"},
		})

		n := client.Project().Get(c.Args()[0])
		n.Credentials = creds
		client.Project().Update(n)
		color.Green("Project " + n.Name + " successfully updated")

		return nil
	},
}
