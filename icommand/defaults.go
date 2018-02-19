package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/urfave/cli"
)

func defaultsICmd(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	useCmd := &ishell.Cmd{
		Name: "default",
		Help: "sets defaults to work with ernest",
	}
	useCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "Lists all defaults on the current client",
		Func: func(c *ishell.Context) {
			for k, v := range defaults {
				c.Println(" - ", k, " : ", v)
			}
		},
	})
	useCmd.AddCmd(&ishell.Cmd{
		Name: "clear",
		Help: "Use specific project on subsequent queries",
		Func: func(c *ishell.Context) {
			defaults = map[string]string{}
			c.Println("All defaults have been reset")
		},
	})
	useCmd.AddCmd(&ishell.Cmd{
		Name: "project",
		Help: "Use specific project on subsequent queries",
		Func: func(c *ishell.Context) {
			setDefaults(c, map[string]string{"project": "Project name : "})
		},
	})

	useCmd.AddCmd(&ishell.Cmd{
		Name: "environment",
		Help: "Use specific environment on subsequent queries",
		Func: func(c *ishell.Context) {
			setDefaults(c, map[string]string{"environment": "Environment name : ", "project": "Project name : "})
		},
	})

	return useCmd
}
