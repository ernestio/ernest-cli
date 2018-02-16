package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func notify(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name: "notify",
		Help: "Notifications management",
	}

	cmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: h.T("notification.list.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.ListNotifications.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: h.T("notification.delete.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			args = mapArgs(c, map[string]input{
				"name": input{out: "Notification name : "},
			})
			command.DeleteNotification.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "update",
		Help: h.T("notification.update.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			args = mapArgs(c, map[string]input{
				"name":   input{out: "Notification name : "},
				"config": input{out: "Configupration file path : "},
			})
			command.UpdateNotification.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: h.T("notification.create.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			args = mapArgs(c, map[string]input{
				"name":   input{out: "Notification name : "},
				"type":   input{out: "Notification type : "},
				"config": input{out: "Configupration file path : "},
			})
			command.CreateNotification.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: h.T("notification.service.add.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			args = mapArgs(c, map[string]input{
				"name":        input{out: "Notification name : "},
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})
			command.CreateNotification.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "remove",
		Help: h.T("notification.service.rm.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			args = mapArgs(c, map[string]input{
				"name":        input{out: "Notification name : "},
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})
			command.CreateNotification.Run(getContext(ctx, args, flags))
		},
	})

	return cmd
}
