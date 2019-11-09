package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func loginICmd(shell *ishell.Shell) func(c *ishell.Context) {
	return func(c *ishell.Context) {
		// disable the '>>>' for cleaner same line input.
		c.ShowPrompt(false)
		defer c.ShowPrompt(true) // yes, revert after login.

		cfg, err := login(c)
		if err != nil {
			c.Println("Invalid credentials")
			return
		}

		userChain = append(userChain, cfg)
		updatePrompt(shell)
	}
}

func logout(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "logout",
		Help: h.T("logout.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.Logout.Run(getContext(ctx, args, flags))
		},
	}
}
