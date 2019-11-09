package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func users(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name: "user",
		Help: "Users management",
	}

	cmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: h.T("user.list.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.ListUsers.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: h.T("user.create.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"email": input{out: "Email : "},
			})
			args = mapArgs(c, map[string]input{
				"username": input{out: "Username : "},
				"pasword":  input{out: "Password : "},
			})
			command.CreateUser.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "password",
		Help: h.T("user.password.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"user":     input{out: "Username : "},
				"password": input{out: "Password : "},
			})
			command.DeletePolicy.Run(getContext(ctx, args, flags))

			command.PasswordUser.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "disable",
		Help: h.T("user.disable.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"user": input{out: "Username : "},
			})
			command.UpdateNotification.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "info",
		Help: h.T("user.info.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.InfoUser.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "add-admin",
		Help: h.T("user.admin.add.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"user": input{out: "Username : "},
			})
			command.AddAdminUser.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "rm-admin",
		Help: h.T("user.admin.rm.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"user": input{out: "Username : "},
			})
			command.RmAdminUser.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "enable-mfa",
		Help: h.T("user.enable-mfa.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"user-name": input{out: "Username : "},
			})
			command.EnableMFA.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "disable-mfa",
		Help: h.T("user.disable-mfa.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"user-name": input{out: "Username : "},
			})
			command.DisableMFA.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "reset-mfa",
		Help: h.T("user.reset-mfa.description"),
		Func: func(c *ishell.Context) {
			var args []string
			flags := mapFlags(c, map[string]input{
				"user-name": input{out: "Username : "},
			})
			command.ResetMFA.Run(getContext(ctx, args, flags))
		},
	})

	return cmd
}
