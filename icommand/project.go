package icommand

import (
	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	"github.com/urfave/cli"
)

var (
	awsType    = "aws"
	azureType  = "aws"
	vcloudType = "aws"
)

func projectICmd(shell *ishell.Shell, ctx *cli.Context) *ishell.Cmd {
	projectCmd := &ishell.Cmd{
		Name: "project",
		Help: "manage ernest projects",
	}

	// ListProjects
	projectCmd.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "List all your projects on ernest",
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			command.ListProjects.Run(getContext(ctx, args, flags))
		},
	})

	// Info
	projectCmd.AddCmd(&ishell.Cmd{
		Name: "info",
		Help: "Get project information",
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			args = mapArgs(c, map[string]input{
				"project": input{out: "Project name : "},
			})
			execute(command.InfoProject, getContext(ctx, args, flags))
		},
	})

	// Update
	projectCmd.AddCmd(&ishell.Cmd{
		Name: "update",
		Help: "Updates an specific project",
		Func: func(c *ishell.Context) {
			var args []string
			var flags map[string]string
			inputArgs := map[string]input{
				"type [aws, azure, vcloud]": input{out: "Project type : "},
			}
			preArgs := mapArgs(c, inputArgs)
			inputArgs = map[string]input{
				"project": input{out: "Project name : "},
			}
			args = mapArgs(c, inputArgs)
			inputFlags := projectInputFlags(preArgs[1])
			flags = mapFlags(c, inputFlags)

			switch preArgs[1] {
			case awsType:
				execute(command.UpdateAWSProject, getContext(ctx, args, flags))
			case azureType:
				execute(command.UpdateAzureProject, getContext(ctx, args, flags))
			case vcloudType:
				execute(command.UpdateVCloudProject, getContext(ctx, args, flags))
			}
		},
	})

	// Create
	projectCmd.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "Creates a project",
		Func: func(c *ishell.Context) {
			projectType, args := projectInputArgs(c)
			inputFlags := projectInputFlags(projectType)
			inputFlags["region"] = input{out: "Region : "}
			flags := mapFlags(c, inputFlags)

			switch projectType {
			case awsType:
				execute(command.CreateAWSProject, getContext(ctx, args, flags))
			case azureType:
				execute(command.CreateAzureProject, getContext(ctx, args, flags))
			case vcloudType:
				execute(command.CreateVcloudProject, getContext(ctx, args, flags))
			}
		},
	})

	return projectCmd
}

func projectInputArgs(c *ishell.Context) (projectType string, args []string) {
	inputArgs := map[string]input{
		"type [aws, azure, vcloud]": input{out: "Project type : "},
	}
	preArgs := mapArgs(c, inputArgs)
	projectType = preArgs[1]
	inputArgs = map[string]input{
		"project": input{out: "Project name : "},
	}
	args = mapArgs(c, inputArgs)
	return
}

func projectInputFlags(projectType string) (inputArgs map[string]input) {
	switch projectType {
	case "aws":
		inputArgs = map[string]input{
			"access_key_id":     input{out: "AWS access key ID : "},
			"secret_access_key": input{out: "AWS secret access key : "},
		}
	case "azure":
		inputArgs = map[string]input{
			"region":                input{out: "Region : "},
			"azure_subscription_id": input{out: "Subscription ID : "},
			"azure_client_id":       input{out: "Client ID : "},
			"azure_client_secret":   input{out: "Client Secret : "},
			"azure_tenant_id":       input{out: "Tenant ID : "},
			"azure_environment":     input{out: "Environment : "},
		}
	case "vcloud":
		inputArgs = map[string]input{
			"vcloud_url": input{out: "URL : "},
			"user":       input{out: "User : "},
			"org":        input{out: "Organisation : "},
			"vdc":        input{out: "VDC : "},
			"password":   input{out: "Password : "},
		}
	}

	return inputArgs
}
