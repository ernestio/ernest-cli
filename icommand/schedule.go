package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func schedule(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name: "schedule",
		Help: "Schedules management",
	}

	cmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: h.T("envs.schedule.list.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})
			command.EnvListSchedules.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: h.T("envs.schedule.rm.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})
			command.EnvRmSchedule.Run(getContext(ctx, args, flags))
		},
	})

	cmd.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: h.T("envs.schedule.add.description"),
		Func: func(c *ishell.Context) {
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
				"schedule":    input{out: "Name : "},
			})
			flags := mapFlags(c, map[string]input{
				"action":        input{out: "Action : "},
				"interval":      input{out: "Schedule : "},
				"instance_type": input{out: "Instance type : "},
			})
			command.EnvAddSchedule.Run(getContext(ctx, args, flags))
		},
	})

	return cmd
}
