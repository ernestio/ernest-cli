/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"fmt"

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
		tStringFlag("aws.create.flags.region"),
		tStringFlag("aws.create.flags.access_key_id"),
		tStringFlag("aws.create.flags.secret_access_key"),
		tStringFlag("aws.create.flags.template"),
		tBoolFlag("aws.create.flags.fake"),
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
		color.Green(fmt.Sprintf(h.T("aws.create.success"), p.Name))

		return nil
	},
}

// UpdateAWSProject : Updates the specified VCloud project
var UpdateAWSProject = cli.Command{
	Name:        "aws",
	Usage:       h.T("aws.update.usage"),
	ArgsUsage:   h.T("aws.update.args"),
	Description: h.T("aws.update.description"),
	Flags: []cli.Flag{
		tStringFlag("aws.update.flags.access_key_id"),
		tStringFlag("aws.update.flags.secret_access_key"),
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
		color.Green(fmt.Sprintf(h.T("aws.create.success"), n.Name))

		return nil
	},
}
