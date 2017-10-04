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

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListEnvs ...
var ListEnvs = cli.Command{
	Name:        "list",
	Usage:       h.T("envs.list.usage"),
	ArgsUsage:   h.T("envs.list.args"),
	Description: h.T("envs.list.description"),
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

		view.PrintEnvList(envs)
		return nil
	},
}

// UpdateEnv :  Creates an empty environment based on a name and a project name
var UpdateEnv = cli.Command{
	Name:        "update",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.update.usage"),
	ArgsUsage:   h.T("envs.update.args"),
	Description: h.T("envs.update.description"),
	Flags:       AllProviderFlags,
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

		err := m.UpdateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c))
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		color.Green("Environment successfully updated")

		return nil
	},
}

// CreateEnv :  Creates an empty environment based on a name and a project name
var CreateEnv = cli.Command{
	Name:        "create",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.create.usage"),
	ArgsUsage:   h.T("envs.create.args"),
	Description: h.T("envs.create.description"),
	Flags: append([]cli.Flag{
		cli.StringFlag{
			Name:  "credentials",
			Usage: "will override project information",
		},
	}, AllProviderFlags...),
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

		err := m.CreateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c))
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Environment successfully created")

		return nil
	},
}

// ApplyEnv command
// Applies changes described on a YAML file to an env
var ApplyEnv = cli.Command{
	Name:        "apply",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.apply.usage"),
	ArgsUsage:   h.T("envs.apply.args"),
	Description: h.T("envs.apply.description"),
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
			os.Exit(1)
		}
		if dry == true {
			fmt.Println(string(response))
		}
		return nil
	},
}

// DestroyEnv command
var DestroyEnv = cli.Command{
	Name:        "delete",
	Aliases:     []string{"destroy", "d"},
	Usage:       h.T("envs.destroy.usage"),
	ArgsUsage:   h.T("envs.destroy.args"),
	Description: h.T("envs.destroy.description"),
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
	Name:        "history",
	Usage:       h.T("envs.history.usage"),
	ArgsUsage:   h.T("envs.history.args"),
	Description: h.T("envs.history.description"),
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
		view.PrintEnvHistory(env, envs)
		return nil
	},
}

// ResetEnv command
var ResetEnv = cli.Command{
	Name:        "reset",
	Usage:       h.T("envs.reset.usage"),
	ArgsUsage:   h.T("envs.reset.args"),
	Description: h.T("envs.reset.description"),
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
	Name:        "revert",
	Usage:       h.T("envs.revert.usage"),
	ArgsUsage:   h.T("envs.revert.args"),
	Description: h.T("envs.revert.description"),
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
	Name:        "definition",
	Aliases:     []string{"s"},
	Usage:       h.T("envs.definition.usage"),
	ArgsUsage:   h.T("envs.definition.args"),
	Description: h.T("envs.definition.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "build",
			Value: "",
			Usage: "Build ID",
		},
	},
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
			definition, err := m.BuildDefinitionFromIndex(cfg.Token, project, env, c.String("build"))
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}
			fmt.Println(string(definition))
		} else {
			definition, err := m.LatestBuildDefinition(cfg.Token, project, env)
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}

			fmt.Println(string(definition))
		}
		return nil
	},
}

// InfoEnv : Shows detailed information of a env by its name
var InfoEnv = cli.Command{
	Name:        "info",
	Aliases:     []string{"i"},
	Usage:       h.T("envs.info.usage"),
	ArgsUsage:   h.T("envs.info.args"),
	Description: h.T("envs.info.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "build",
			Value: "",
			Usage: "Build ID",
		},
	},
	Action: func(c *cli.Context) error {
		var err error
		var b model.Build

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
			b, err = m.BuildStatus(cfg.Token, project, env, build)
		} else {
			b, err = m.LatestBuildStatus(cfg.Token, project, env)
		}

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		view.PrintEnvInfo(&b)
		return nil
	},
}

// DiffEnv : Shows detailed information of an env by its name
var DiffEnv = cli.Command{
	Name:        "diff",
	Aliases:     []string{"i"},
	Usage:       h.T("envs.diff.usage"),
	ArgsUsage:   h.T("envs.diff.args"),
	Description: h.T("envs.diff.description"),
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

		project := c.Args()[0]
		env := c.Args()[1]
		b1 := c.Args()[2]
		b2 := c.Args()[3]

		build1, err := m.BuildStatus(cfg.Token, project, env, b1)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		build2, err := m.BuildStatus(cfg.Token, project, env, b2)
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
	Name:        "import",
	Aliases:     []string{"i"},
	Usage:       h.T("envs.import.usage"),
	ArgsUsage:   h.T("envs.import.args"),
	Description: h.T("envs.import.description"),
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
