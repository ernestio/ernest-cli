package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func logger(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name: "logger",
		Help: "Logger management",
	}

	cmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: h.T("logger.list.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.ListLoggers.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: h.T("logger.set.description"),
		Func: func(c *ishell.Context) {
			args := mapArgs(c, map[string]input{
				"name": input{out: "Logger name : "},
			})
			flags := mapFlags(c, map[string]input{
				"logfile":  input{out: "Log File : "},
				"hostname": input{out: "Hostname : "},
				"port":     input{out: "Port : "},
				"timeout":  input{out: "Timeout : "},
				"token":    input{out: "Token : "},
				"env":      input{out: "Environment : "},
			})
			command.SetLogger.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: h.T("logger.del.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"name": input{out: "Notification name : "},
			})
			command.DelLogger.Run(getContext(ctx, args, flags))
		},
	})

	return cmd
}
