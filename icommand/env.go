package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/urfave/cli"
)

func envICmd(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	envCmd := &ishell.Cmd{
		Name: "env",
		Help: "manage ernest environments",
	}

	// List
	envCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: h.T("envs.list.description"),
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.ListEnvs.Run(getContext(ctx, args, flags))
		},
	})

	// Create
	envCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: h.T("envs.create.description"),
		Func: func(c *ishell.Context) {
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})
			flags := mapEnvironmentFlags(c)
			command.CreateEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Update
	envCmd.AddCmd(&ishell.Cmd{
		Name: "update",
		Help: h.T("envs.update.description"),
		Func: func(c *ishell.Context) {
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})
			flags := mapEnvironmentFlags(c)
			command.UpdateEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Apply
	envCmd.AddCmd(&ishell.Cmd{
		Name: "apply",
		Help: h.T("envs.create.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"template": input{out: "Template : "},
			})

			command.ApplyEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Sync
	envCmd.AddCmd(&ishell.Cmd{
		Name: "sync",
		Help: h.T("envs.sync.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.SyncEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Review
	envCmd.AddCmd(&ishell.Cmd{
		Name: "review",
		Help: h.T("envs.review.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.ReviewEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Resolve
	envCmd.AddCmd(&ishell.Cmd{
		Name: "resolve",
		Help: h.T("envs.resolve.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.ResolveEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Delete
	envCmd.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: h.T("envs.delete.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.DestroyEnv.Run(getContext(ctx, args, flags))
		},
	})

	// History
	envCmd.AddCmd(&ishell.Cmd{
		Name: "history",
		Help: h.T("envs.history.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.HistoryEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Reset
	envCmd.AddCmd(&ishell.Cmd{
		Name: "reset",
		Help: h.T("envs.reset.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.ResetEnv.Run(getContext(ctx, args, flags))
		},
	})

	// Revert
	envCmd.AddCmd(&ishell.Cmd{
		Name: "revert",
		Help: h.T("envs.revert.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
				"build":       input{out: "Build : "},
			})

			command.RevertEnv.Run(getContext(ctx, args, flags))
		},
	})

	// definition
	envCmd.AddCmd(&ishell.Cmd{
		Name: "definition",
		Help: h.T("envs.definition.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.DefinitionEnv.Run(getContext(ctx, args, flags))
		},
	})

	// info
	envCmd.AddCmd(&ishell.Cmd{
		Name: "info",
		Help: h.T("envs.info.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.InfoEnv.Run(getContext(ctx, args, flags))
		},
	})

	// import
	envCmd.AddCmd(&ishell.Cmd{
		Name: "import",
		Help: h.T("envs.import.description"),
		Func: func(c *ishell.Context) {
			var flags map[string]string
			args := mapArgs(c, map[string]input{
				"project":     input{out: "Project : "},
				"environment": input{out: "Environment : "},
			})

			command.ImportEnv.Run(getContext(ctx, args, flags))
		},
	})

	return envCmd
}

func mapEnvironmentFlags(c *ishell.Context) map[string]string {
	preArgs := mapArgs(c, map[string]input{
		"type [aws, azure, vcloud]": input{out: "What's the type of the project (aws, azure, vcloud) : "},
	})

	inputArgs := map[string]input{
		"credentials":   input{out: "Credentials : "},
		"sync_interval": input{out: "Sync Interval : "},
		"submissions":   input{out: "Submissions : "},
	}
	switch preArgs[1] {
	case awsType:
		inputArgs["region"] = input{out: "Project region"}
		inputArgs["access_key_id"] = input{out: "AWS access key id : "}
		inputArgs["secret_access_key"] = input{out: "AWS Secret access key : "}
	case vcloudType:
		inputArgs["user"] = input{out: "VCloud User : "}
		inputArgs["password"] = input{out: "Your VCloud valid password : "}
		inputArgs["org"] = input{out: "Your vCloud Organization : "}
		inputArgs["vse-url"] = input{out: "VSE URL : "}
		inputArgs["vcloud-url"] = input{out: "VCloud URL"}
		inputArgs["public-network"] = input{out: "Public Network : "}
		inputArgs["vcloud-region"] = input{out: "VCloud region : "}
	case azureType:
		inputArgs["subscription_id"] = input{out: "Azure subscription id"}
		inputArgs["client_id"] = input{out: "Azure client id : "}
		inputArgs["client_secret"] = input{out: "Azure client secret : "}
		inputArgs["tenant_id"] = input{out: "Azure tenant_id : "}
		inputArgs["environment"] = input{out: "Azure environment : "}
	}
	flags := mapFlags(c, inputArgs)
	return flags
}
