package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/urfave/cli"
)

func sudoICmd(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "sudo",
		Help: "Execute a set of commands as another user",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			cfg, err := login(c)
			if err != nil {
				c.Println("Invalid credentials")
				return
			}

			userChain = append(userChain, cfg)
			updatePrompt(shell)
		},
	}
}
