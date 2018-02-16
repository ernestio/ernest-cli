package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func role(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name: "roles",
		Help: "Roles management",
	}

	cmd.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: h.T("roles.set.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"project":     input{out: "Project (empty if not applicable) : "},
				"environment": input{out: "Environment (empty if not applicable) : "},
				"policy":      input{out: "Policy (empty if not applicable) : "},
			})
			command.ListPolicies.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "unset",
		Help: h.T("roles.unset.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"project":     input{out: "Project (empty if not applicable) : "},
				"environment": input{out: "Environment (empty if not applicable) : "},
				"policy":      input{out: "Policy (empty if not applicable) : "},
			})
			command.DeletePolicy.Run(getContext(ctx, args, flags))
		},
	})

	return cmd
}
