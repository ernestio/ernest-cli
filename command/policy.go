/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdPolicy subcommand
import (
	"fmt"
	"io/ioutil"
	"strings"

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

		client.Policy().CreateDocument(name, string(spec))
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
			Name: flags["policy-name"].(string),
		}
		client.Policy().Create(&policy)
		client.Policy().CreateDocument(policy.Name, string(spec))
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
		tStringFlag("policy.show.flags.revision"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
			"revision":    flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)

		n := client.Policy().GetDocument(flags["policy-name"].(string), flags["revision"].(string))
		fmt.Println(n.Definition)
		return nil
	},
}

// HistoryPolicy : Display a policys revisions
var HistoryPolicy = cli.Command{
	Name:        "history",
	Usage:       h.T("policy.history.usage"),
	ArgsUsage:   h.T("policy.history.args"),
	Description: h.T("policy.history.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.history.flags.name"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)

		documents := client.Policy().ListDocuments(flags["policy-name"].(string))
		view.PrintPolicyHistory(documents)

		return nil
	},
}

// AttachPolicy : Display an existing policy
var AttachPolicy = cli.Command{
	Name:        "attach",
	Usage:       h.T("policy.attach.usage"),
	ArgsUsage:   h.T("policy.attach.args"),
	Description: h.T("policy.attach.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.attach.flags.name"),
		tStringFlag("policy.attach.flags.environment"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
			"environment": flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)
		env := flags["environment"].(string)
		parts := strings.Split(env, "/")
		if len(parts) != 2 {
			h.PrintError(h.T("policy.attach.errors.invalid_name"))
		}

		p := client.Policy().Get(flags["policy-name"].(string))
		_ = client.Environment().Get(parts[0], parts[1])
		for _, v := range p.Environments {
			if v == env {
				h.PrintError(h.T("policy.attach.errors.already_attached"))
			}
		}
		p.Environments = append(p.Environments, env)
		client.Policy().Update(p)

		color.Green(fmt.Sprintf(h.T("policy.attach.success"), p.Name, env))
		return nil
	},
}

// DetachPolicy : Display an existing policy
var DetachPolicy = cli.Command{
	Name:        "detach",
	Usage:       h.T("policy.detach.usage"),
	ArgsUsage:   h.T("policy.detach.args"),
	Description: h.T("policy.detach.description"),
	Flags: []cli.Flag{
		tStringFlag("policy.detach.flags.name"),
		tStringFlag("policy.detach.flags.environment"),
	},
	Action: func(c *cli.Context) error {
		flags := parseTemplateFlags(c, map[string]flagDef{
			"policy-name": flagDef{typ: "string", req: true},
			"environment": flagDef{typ: "string", req: true},
		})
		client := esetup(c, AuthUsersValidation)
		env := flags["environment"].(string)
		parts := strings.Split(env, "/")
		if len(parts) != 2 {
			h.PrintError(h.T("policy.detach.errors.invalid_name"))
		}

		p := client.Policy().Get(flags["policy-name"].(string))
		_ = client.Environment().Get(parts[0], parts[1])
		var toBeAttached []string
		for _, v := range p.Environments {
			if v != env {
				toBeAttached = append(toBeAttached, v)
			}
		}
		if len(toBeAttached) == len(p.Environments) {
			h.T("policy.detach.error.not_attached")
		}
		p.Environments = toBeAttached
		client.Policy().Update(p)

		color.Green(fmt.Sprintf(h.T("policy.detach.success"), p.Name, env))
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
		HistoryPolicy,
		AttachPolicy,
		DetachPolicy,
	},
}
