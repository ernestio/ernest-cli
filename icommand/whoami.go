package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/urfave/cli"
)

func whoamiICmd(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "whoami",
		Help: "display effective user id",
		Func: func(c *ishell.Context) {
			c.Println(userChain[len(userChain)-1].User)
		},
	}
}
