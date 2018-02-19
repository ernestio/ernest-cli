package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func policy(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name: "policy",
		Help: "Policy management",
	}

	cmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: h.T("policy.list.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.ListPolicies.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: h.T("policy.delete.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"policy-name": input{out: "Policy name : "},
			})
			command.DeletePolicy.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "update",
		Help: h.T("policy.update.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"policy-name": input{out: "Policy name : "},
				"spec":        input{out: "Spec : "},
			})
			command.UpdatePolicy.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: h.T("policy.update.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"policy-name": input{out: "Policy name : "},
			})
			command.UpdatePolicy.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: h.T("policy.show.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"policy-name": input{out: "Policy name : "},
				"spec":        input{out: "Spec : "},
			})
			command.CreatePolicy.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "attach",
		Help: h.T("policy.attach.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"name":        input{out: "Policy name : "},
				"environment": input{out: "Environment : "},
			})
			command.AttachPolicy.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "detach",
		Help: h.T("policy.detach.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"name":        input{out: "Policy name : "},
				"environment": input{out: "Environment : "},
			})
			command.AttachPolicy.Run(getContext(ctx, args, flags))
		},
	})

	return cmd
}
