package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	"github.com/urfave/cli"
)

func info(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "info",
		Help: "Info for the current session",
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.Info.Run(getContext(ctx, args, flags))
		},
	}
}
