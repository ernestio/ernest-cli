package icommand

import (
	"github.com/abiosoft/ishell"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/urfave/cli"
)

func exitICmd(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "exit",
		Help: "Exit current session",
		Func: func(c *ishell.Context) {
			if len(userChain) <= 1 {
				c.Println("Bye!")
				c.Stop()
			} else {
				c.Println("Bye", userChain[len(userChain)-1].User, "!")
				i := len(userChain) - 1
				userChain = append(userChain[:i], userChain[i+1:]...)
				if err := model.SaveConfig(userChain[len(userChain)-1]); err != nil {
					h.PrintError("Can't write config file")
				}
				updatePrompt(shell)
			}
		},
	}
}
