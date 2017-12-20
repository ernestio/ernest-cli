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

// CreateAzureProject : Creates an AWS project
var CreateAzureProject = cli.Command{
	Name:        "azure",
	Usage:       h.T("azure.create.usage"),
	Description: h.T("azure.create.description"),
	ArgsUsage:   h.T("azure.create.args"),
	Flags: []cli.Flag{
		tStringFlag("azure.create.flags.region"),
		tStringFlag("azure.create.flags.subscription_id"),
		tStringFlag("azure.create.flags.client_id"),
		tStringFlag("azure.create.flags.client_secret"),
		tStringFlag("azure.create.flags.tenant_id"),
		tStringFlag("azure.create.flags.environment"),
		tBoolFlag("azure.create.flags.fake"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "azure.create.args")
		client := esetup(c, AuthUsersValidation)
		creds := parseTemplateFlags(c, map[string]flagDef{
			"region":          flagDef{typ: "string", mapto: "region", req: true},
			"subscription_id": flagDef{typ: "string", mapto: "azure_subscription_id", req: true},
			"client_id":       flagDef{typ: "string", mapto: "azure_client_id", req: true},
			"client_secret":   flagDef{typ: "string", mapto: "azure_client_secret", req: true},
			"tenant_id":       flagDef{typ: "string", mapto: "azure_tenant_id", req: true},
			"environment":     flagDef{typ: "string", mapto: "azure_environment", def: "default", req: true},
			"fake":            flagDef{typ: "bool", def: false},
		})
		rtype := "azure"
		if creds["fake"].(bool) {
			rtype = "azure-fake"
		}

		p := &emodels.Project{
			Name:        c.Args()[0],
			Type:        rtype,
			Credentials: creds,
		}
		client.Project().Create(p)
		color.Green(fmt.Sprintf(h.T("azure.create.success"), p.Name))

		return nil
	},
}

// UpdateAzureProject : Updates the specified VCloud project
var UpdateAzureProject = cli.Command{
	Name:        "azure",
	Usage:       h.T("azure.update.usage"),
	Description: h.T("azure.update.description"),
	ArgsUsage:   h.T("azure.update.args"),
	Flags: []cli.Flag{
		tStringFlag("azure.update.flags.subscription_id"),
		tStringFlag("azure.update.flags.client_id"),
		tStringFlag("azure.update.flags.client_secret"),
		tStringFlag("azure.update.flags.tenant_id"),
		tStringFlag("azure.update.flags.environment"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "azure.update.args")
		client := esetup(c, AuthUsersValidation)

		creds := parseTemplateFlags(c, map[string]flagDef{
			"subscription_id": flagDef{typ: "string", mapto: "azure_subscription_id"},
			"client_id":       flagDef{typ: "string", mapto: "azure_client_id"},
			"client_secret":   flagDef{typ: "string", mapto: "azure_client_secret"},
			"tenant_id":       flagDef{typ: "string", mapto: "azure_tenant_id"},
			"environment":     flagDef{typ: "string", mapto: "azure_environment"},
			"region":          flagDef{typ: "string", mapto: "region"},
		})

		n := client.Project().Get(c.Args()[0])
		n.Credentials = creds
		client.Project().Update(n)
		color.Green(fmt.Sprintf(h.T("azure.update.success"), n.Name))

		return nil
	},
}
