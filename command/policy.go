/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdPolicy subcommand
import (
	"fmt"
	"io/ioutil"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// ListPolicies ...
var ListPolicies = cli.Command{
	Name:        "list",
	Usage:       h.T("policy.list.usage"),
	ArgsUsage:   h.T("policy.list.args"),
	Description: h.T("policy.list.description"),
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		policies := client.Policy().List()

		view.PrintPolicyList(policies)

		return nil
	},
}

// DeletePolicy : Will delete the specified policy
var DeletePolicy = cli.Command{
	Name:        "delete",
	Usage:       h.T("policy.delete.usage"),
	ArgsUsage:   h.T("policy.delete.args"),
	Description: h.T("policy.delete.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.delete.flags.name"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
		})

		client := esetup(c, AuthUsersValidation)
		name := flags["policy-name"].(string)

		client.Policy().Delete(name)
		color.Green(fmt.Sprintf(h.T("policy.delete.success"), name))
		return nil
	},
}

// UpdatePolicy : Will update the policy specific fields
var UpdatePolicy = cli.Command{
	Name:        "update",
	Usage:       h.T("policy.update.usage"),
	ArgsUsage:   h.T("policy.update.args"),
	Description: h.T("policy.update.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.update.flags.name"),
		tStringFlag("policy.update.flags.spec"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
			"spec":        flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)
		name := flags["policy-name"].(string)
		spec, err := ioutil.ReadFile(flags["spec"].(string))
		if err != nil {
			h.PrintError(h.T("policy.update.errors.spec"))
		}

		n := client.Policy().Get(name)
		n.Definition = string(spec)
		client.Policy().Update(n)
		color.Green(fmt.Sprintf(h.T("policy.update.success"), name))
		return nil
	},
}

// CreatePolicy : Creates a new user
var CreatePolicy = cli.Command{
	Name:        "create",
	Usage:       h.T("policy.create.usage"),
	ArgsUsage:   h.T("policy.create.args"),
	Description: h.T("policy.create.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.create.flags.name"),
		tStringFlag("policy.create.flags.spec"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
			"spec":        flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)
		spec, err := ioutil.ReadFile(flags["spec"].(string))
		if err != nil {
			h.PrintError(h.T("policy.create.errors.spec"))
		}

		policy := emodels.Policy{
			Name:       flags["policy-name"].(string),
			Definition: string(spec),
		}
		client.Policy().Create(&policy)
		color.Green(fmt.Sprintf(h.T("policy.create.success"), policy.Name))
		return nil
	},
}

// ShowPolicy : Display an existing policy
var ShowPolicy = cli.Command{
	Name:        "show",
	Usage:       h.T("policy.show.usage"),
	ArgsUsage:   h.T("policy.show.args"),
	Description: h.T("policy.show.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.show.flags.name"),
		tStringFlag("policy.show.flags.spec"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)

		n := client.Policy().Get(flags["policy-name"].(string))
		fmt.Println(n.Definition)
		return nil
	},
}

// CmdPolicy ...
var CmdPolicy = cli.Command{
	Name:    "policy",
	Aliases: []string{"p"},
	Usage:   "Policy related subcommands",
	Subcommands: []cli.Command{
		ListPolicies,
		CreatePolicy,
		UpdatePolicy,
		DeletePolicy,
		ShowPolicy,
	},
}
