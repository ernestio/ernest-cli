/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListEnvs ...
var ListEnvs = cli.Command{
	Name:      "list",
	Usage:     "List available environments.",
	ArgsUsage: " ",
	Description: `List available environments and shows its most relevant information.

   Example:
    $ ernest environment list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		envs, err := m.ListEnvs(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		view.PrintEnvsList(envs)
		return nil
	},
}

// UpdateEnv :  Creates an empty environment based on a name and a project name
var UpdateEnv = cli.Command{
	Name:      "update",
	Aliases:   []string{"a"},
	ArgsUsage: "<project> <environment>",
	Usage:     "Creates an empty environment based on a specific project",
	Flags:     AllProviderFlags,
	Description: `You must be logged in to execute this command.

   Examples:
    $ ernest env update --credentials project.yml my_project my_environment
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You must provide the project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You must provide the new environment name")
			return nil
		}
		project := c.Args()[0]
		env := c.Args()[1]

		_, err := m.UpdateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c))
		if err != nil {
			color.Red(err.Error())
		}
		fmt.Println("Environment successfully updated")

		return nil
	},
}

// CreateEnv :  Creates an empty environment based on a name and a project name
var CreateEnv = cli.Command{
	Name:      "create",
	Aliases:   []string{"a"},
	ArgsUsage: "<project> <environment>",
	Usage:     "Creates an empty environment based on a specific project",
	Flags: append([]cli.Flag{
		cli.StringFlag{
			Name:  "credentials",
			Usage: "will override project information",
		},
	}, AllProviderFlags...),
	Description: `You must be logged in to execute this command.

   Examples:
    $ ernest env create my_project my_environment
    $ ernest env create --credentials project.yml my_project my_environment
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You must provide the project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You must provide the new environment name")
			return nil
		}
		project := c.Args()[0]
		env := c.Args()[1]

		_, err := m.CreateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c))
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		fmt.Println("Environment successfully created")

		return nil
	},
}

// ApplyEnv command
// Applies changes described on a YAML file to an env
var ApplyEnv = cli.Command{
	Name:      "apply",
	Aliases:   []string{"a"},
	ArgsUsage: "<file.yml>",
	Usage:     "Builds or changes infrastructure.",
	Flags: append([]cli.Flag{
		cli.BoolFlag{
			Name:  "dry",
			Usage: "print the changes to be applied on an environment intead of applying them",
		},
		cli.StringFlag{
			Name:  "credentials",
			Usage: "will override project information",
		},
	}, AllProviderFlags...),
	Description: `Sends an environment YAML description file to Ernest to be executed.
   You must be logged in to execute this command.

   If the file is not provided, ernest.yml will be used by default.

   Examples:
    $ ernest env apply myenvironment.yml
    $ ernest env apply --dry myenvironment.yml
	`,
	Action: func(c *cli.Context) error {
		file := "ernest.yml"
		if len(c.Args()) == 1 {
			file = c.Args()[0]
		}
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		var err error
		dry := c.Bool("dry")
		monit := true
		if dry == true {
			monit = false
		}
		response, err := m.Apply(cfg.Token, file, ProviderFlagsToSlice(c), monit, dry)
		if err != nil {
			color.Red(err.Error())
		}
		if dry == true {
			fmt.Println(string(response))
		}
		return nil
	},
}

// DestroyEnv command
var DestroyEnv = cli.Command{
	Name:      "delete",
	Aliases:   []string{"destroy", "d"},
	ArgsUsage: "<project> <environment_name>",
	Usage:     "Destroy an environment.",
	Description: `Destroys an environment by name.

   Example:
    $ ernest env delete <my_project> <my_environment>
  `,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force,f",
			Usage: "Hard ernest env removal.",
		},
		cli.BoolFlag{
			Name:  "yes,y",
			Usage: "Destroy an environment without prompting confirmation.",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify an existing project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify an existing project environment")
			return nil
		}
		project := c.Args()[0]
		env := c.Args()[1]

		if c.Bool("force") {
			err := m.ForceDestroy(cfg.Token, project, env)
			if err != nil {
				color.Red(err.Error())
				return nil
			}
		} else {
			if c.Bool("yes") {
				err := m.Destroy(cfg.Token, project, env, true)
				if err != nil {
					color.Red(err.Error())
					return nil
				}
			} else {
				fmt.Print("Do you really want to destroy this environment? (Y/n) ")
				if askForConfirmation() == false {
					return nil
				}
				err := m.Destroy(cfg.Token, project, env, true)
				if err != nil {
					color.Red(err.Error())
					return nil
				}
			}
		}
		color.Green("Environment successfully removed")
		return nil
	},
}

// HistoryEnv command
// Shows the history of an env, a list of builds
var HistoryEnv = cli.Command{
	Name:      "history",
	Usage:     "Shows the history of an environment, a list of builds",
	ArgsUsage: "ernest-cli env history <my_project> <my_env>",
	Description: `Shows the history of an environment, a list of builds and its status and basic information.

   Example:
    $ ernest env history <my_project> <my_env>
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify an existing project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify an existing environment name")
			return nil
		}

		project := c.Args()[0]
		env := c.Args()[1]

		envs, _ := m.ListBuilds(project, env, cfg.Token)
		view.PrintEnvHistory(envs)
		return nil
	},
}

// ResetEnv command
var ResetEnv = cli.Command{
	Name:      "reset",
	ArgsUsage: "<env_name>",
	Usage:     "Reset an in progress environment.",
	Description: `Reseting an environment creation may cause problems, please make sure you know what are you doing.

   Example:
    $ ernest env reset <my_env>
  `,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify the project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify the environment name")
			return nil
		}
		project := c.Args()[0]
		env := c.Args()[1]
		err := m.ResetEnv(project, env, cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Red("You've successfully resetted the environment '" + project + " / " + env + "'")

		return nil
	},
}

// RevertEnv command
var RevertEnv = cli.Command{
	Name:      "revert",
	ArgsUsage: "<project> <env_name> <build_id>",
	Usage:     "Reverts an environment to a previous state",
	Description: `Reverts an environment to a previous known state using a build ID from 'ernest service history'.

   Example:
    $ ernest env revert <project> <env_name> <build_id>
    $ ernest env revert --dry <project> <env_name> <build_id>
  `,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry",
			Usage: "print the changes to be applied on an environment intead of applying them",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 3 {
			color.Red("Please specify a project, environment and build ID")
			return nil
		}
		project := c.Args()[0]
		env := c.Args()[1]
		buildID := c.Args()[2]
		dry := c.Bool("dry")

		response, err := m.RevertEnv(project, env, buildID, cfg.Token, dry)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		if dry == true {
			fmt.Println(string(response))
		}

		return nil
	},
}

// DefinitionEnv command
// Shows the current definition of an environment by its name
var DefinitionEnv = cli.Command{
	Name:      "definition",
	Aliases:   []string{"s"},
	ArgsUsage: "<project_name> <env_name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "build",
			Value: "",
			Usage: "Build ID",
		},
	},
	Usage: "Show the current definition of an environment by its name",
	Description: `Show the current definition of an environment by its name getting the definition about the build.

   Example:
    $ ernest env definition <my_project> <my_env>
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify the project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify the env name")
			return nil
		}
		project := c.Args()[0]
		env := c.Args()[1]
		if c.String("build") != "" {
			env, err := m.EnvBuildStatus(cfg.Token, project, env, c.String("build"))
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}
			fmt.Println(env.Definition)
		} else {
			env, err := m.EnvStatus(cfg.Token, project, env)
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}

			fmt.Println(env.Definition)
		}
		return nil
	},
}

// InfoEnv : Shows detailed information of a env by its name
var InfoEnv = cli.Command{
	Name:      "info",
	Aliases:   []string{"i"},
	ArgsUsage: "<project_name> <env_name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "build",
			Value: "",
			Usage: "Build ID",
		},
	},
	Usage: "$ ernest env info <my_env> --build <specific build>",
	Description: `Will show detailed information of the last build of a specified service.
	In case you specify --build option you will be able to output the detailed information of specific build of a service.

   Examples:
    $ ernest env definition <my_project> <my_env>
    $ ernest env definition <my_project> <my_env> --build build1
	`,
	Action: func(c *cli.Context) error {
		var err error
		var s model.Service

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify an existing project name")
			return nil
		}
		if len(c.Args()) == 1 {
			color.Red("You should specify an existing env name")
			return nil
		}

		project := c.Args()[0]
		env := c.Args()[1]
		if c.String("build") != "" {
			build := c.String("build")
			s, err = m.EnvBuildStatus(cfg.Token, project, env, build)
		} else {
			s, err = m.EnvStatus(cfg.Token, project, env)
		}

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		view.PrintEnvInfo(&s)
		return nil
	},
}

// DiffEnv : Shows detailed information of an env by its name
var DiffEnv = cli.Command{
	Name:      "diff",
	Aliases:   []string{"i"},
	ArgsUsage: "<env_aname> <build_a> <build_b>",
	Usage:     "$ ernest env diff <project_name> <env_name> <build_a> <build_b>",
	Description: `Will display the diff between two different builds

   Examples:
    $ ernest env diff <my_project> <my_env> 1 2
	`,
	Action: func(c *cli.Context) error {
		var err error

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 4 {
			color.Red("You should specify the project and env names and two build ids to compare them")
			return nil
		}

		env := c.Args()[0]
		project := c.Args()[1]
		b1 := c.Args()[2]
		b2 := c.Args()[3]

		build1, err := m.EnvBuildStatus(cfg.Token, project, env, b1)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		build2, err := m.EnvBuildStatus(cfg.Token, project, env, b2)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		view.PrintEnvDiff(build1, build2)
		return nil
	},
}

// ImportEnv : Shows detailed information of an env by its name
var ImportEnv = cli.Command{
	Name:      "import",
	Aliases:   []string{"i"},
	ArgsUsage: "<env_name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "project",
			Value: "",
			Usage: "Project name",
		},
		cli.StringFlag{
			Name:  "filters",
			Value: "",
			Usage: "Import filters comma delimited list",
		},
	},
	Usage: "$ ernest env import <my_project> <my_env>",
	Description: `Will import the environment <my_env> from project <project_name>

   Examples:
    $ ernest env import my_project my_env
	`,
	Action: func(c *cli.Context) error {
		var err error
		var filters []string

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify an existing project name")
			return nil
		}
		if len(c.Args()) == 1 {
			color.Red("You should specify a valid environment name")
			return nil
		}

		if c.String("filters") != "" {
			filters = strings.Split(c.String("filters"), ",")
		}

		project := c.Args()[0]
		name := c.Args()[1]
		_, err = m.Import(cfg.Token, name, project, filters)

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		return nil
	},
}

func getEnvUUID(output []byte) (string, error) {
	var env struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(output, &env)
	if err != nil {
		return "", err
	}
	return env.ID, nil
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

// CmdEnv ...
var CmdEnv = cli.Command{
	Name:    "environment",
	Aliases: []string{"env"},
	Usage:   "Environment related subcommands",
	Subcommands: []cli.Command{
		ListEnvs,
		CreateEnv,
		UpdateEnv,
		ApplyEnv,
		DestroyEnv,
		HistoryEnv,
		ResetEnv,
		RevertEnv,
		DefinitionEnv,
		InfoEnv,
		MonitorEnv,
		DiffEnv,
		ImportEnv,
	},
}
