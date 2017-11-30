/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"encoding/json"
	"fmt"
	"strconv"
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}
		envs, err := m.ListEnvs(cfg.Token)
		if err != nil {
			h.PrintError(err.Error())
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
	Flags: append([]cli.Flag{
		cli.StringFlag{
			Name:  "sync_interval",
			Usage: "sets the automatic sync interval. Accepts cron syntax, i.e. '@every 1d', '@weekly' or '0 0 * * * *' (Daily at midnight)",
		},
		cli.StringFlag{
			Name:  "submissions",
			Usage: "allows user build submissions from users that have only read only permission to an environment. Options are 'enable' or 'disable'",
		},
	}, AllProviderFlags...),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You must provide the project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You must provide the new environment name")
		}
		project := c.Args()[0]
		env := c.Args()[1]

		e, err := m.EnvStatus(cfg.Token, project, env)
		if err != nil {
			h.PrintError("Environment does not exist!")
		}

		err = m.UpdateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c), MapEnvOptions(c, e.Options), e.Schedules)
		if err != nil {
			h.PrintError(err.Error())
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
		cli.StringFlag{
			Name:  "sync_interval",
			Usage: "sets the automatic sync interval. Accepts cron syntax, i.e. '@every 1d', '@weekly' or '0 0 * * * *' (Daily at midnight)",
		},
		cli.StringFlag{
			Name:  "submissions",
			Usage: "allows user build submissions from users that have only read only permission to an environment. Options are 'enable' or 'disable'",
		},
	}, AllProviderFlags...),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You must provide the project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You must provide the new environment name")
		}
		project := c.Args()[0]
		env := c.Args()[1]

		err := m.CreateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c), MapEnvOptions(c, nil))
		if err != nil {
			h.PrintError(err.Error())
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		var err error
		dry := c.Bool("dry")
		monit := true
		if dry == true {
			monit = false
		}
		response, err := m.Apply(cfg.Token, file, ProviderFlagsToSlice(c), monit, dry)
		if err != nil {
			h.PrintError(err.Error())
		}
		if dry == true {
			fmt.Println(string(response))
		}
		return nil
	},
}

// SyncEnv command
// Syncs an environment with the cloud provider
var SyncEnv = cli.Command{
	Name:        "sync",
	Aliases:     []string{"s"},
	Usage:       h.T("envs.sync.usage"),
	ArgsUsage:   h.T("envs.sync.args"),
	Description: h.T("envs.sync.description"),
	Flags:       append([]cli.Flag{}),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify an existing project environment")
		}

		project := c.Args()[0]
		env := c.Args()[1]

		err := m.SyncEnv(cfg.Token, env, project)
		if err != nil {
			h.PrintError(err.Error())
			return nil
		}

		return nil
	},
}

// ReviewEnv command
// Approval for outstanding build submissions
var ReviewEnv = cli.Command{
	Name:        "review",
	Aliases:     []string{"rev"},
	Usage:       h.T("envs.review.usage"),
	ArgsUsage:   h.T("envs.review.args"),
	Description: h.T("envs.review.description"),
	Flags: append([]cli.Flag{
		cli.BoolFlag{
			Name:  "accept, a",
			Usage: "Accept Sync changes",
		},
		cli.BoolFlag{
			Name:  "reject, r",
			Usage: "Reject Sync changes",
		},
	}),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify an existing project environment")
		}

		project := c.Args()[0]
		env := c.Args()[1]

		var resolution string

		if c.Bool("accept") {
			resolution = "submission-accepted"
		}

		if c.Bool("reject") {
			resolution = "submission-rejected"
		}

		if resolution == "" {
			envs, err := m.ListBuilds(project, env, cfg.Token)
			if err != nil {
				h.PrintError(err.Error())
			}

			id1, err := m.BuildIDFromIndex(cfg.Token, project, env, strconv.Itoa(len(envs)-1))
			if err != nil {
				h.PrintError(err.Error())
			}

			id2, err := m.BuildIDFromIndex(cfg.Token, project, env, strconv.Itoa(len(envs)))
			if err != nil {
				h.PrintError(err.Error())
			}

			def1, err := m.BuildDefinitionByID(cfg.Token, project, env, id1)
			if err != nil {
				h.PrintError(err.Error())
			}

			def2, err := m.BuildDefinitionByID(cfg.Token, project, env, id2)
			if err != nil {
				h.PrintError(err.Error())
			}

			view.PrintEnvDiff(id1, id2, def1, def2)

			return nil
		}

		err := m.ReviewBuild(cfg.Token, env, project, resolution)
		if err != nil {
			h.PrintError(err.Error())
			return nil
		}

		return nil
	},
}

// ResolveEnv command
// Resolves an issue that requires user input
var ResolveEnv = cli.Command{
	Name:        "resolve",
	Aliases:     []string{"r"},
	Usage:       h.T("envs.resolve.usage"),
	ArgsUsage:   h.T("envs.resolve.args"),
	Description: h.T("envs.resolve.description"),
	Flags: append([]cli.Flag{
		cli.BoolFlag{
			Name:  "accept, a",
			Usage: "Accept Sync changes",
		},
		cli.BoolFlag{
			Name:  "reject, r",
			Usage: "Reject Sync changes",
		},
		cli.BoolFlag{
			Name:  "ignore, i",
			Usage: "Ignore Sync changes",
		},
	}),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify an existing project environment")
		}

		var resolution string

		if c.Bool("accept") {
			resolution = "accept-changes"
		}

		if c.Bool("reject") {
			resolution = "reject-changes"
		}

		if c.Bool("ignore") {
			resolution = "ignore-changes"
		}

		if resolution == "" {
			h.PrintError("You should specify a valid resolution [accept|reject|ignore]")
		}

		project := c.Args()[0]
		env := c.Args()[1]

		err := m.ResolveEnv(cfg.Token, env, project, resolution)
		if err != nil {
			h.PrintError(err.Error())
			return nil
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify an existing project environment")
		}
		project := c.Args()[0]
		env := c.Args()[1]

		if c.Bool("force") {
			err := m.ForceDestroy(cfg.Token, project, env)
			if err != nil {
				h.PrintError(err.Error())
			}
		} else {
			if c.Bool("yes") {
				err := m.Destroy(cfg.Token, project, env, true)
				if err != nil {
					h.PrintError(err.Error())
				}
			} else {
				fmt.Print("Do you really want to destroy this environment? (Y/n) ")
				if askForConfirmation() == false {
					return nil
				}
				err := m.Destroy(cfg.Token, project, env, true)
				if err != nil {
					h.PrintError(err.Error())
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify an existing environment name")
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify the project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify the environment name")
		}
		project := c.Args()[0]
		env := c.Args()[1]
		err := m.ResetEnv(project, env, cfg.Token)
		if err != nil {
			h.PrintError(err.Error())
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 3 {
			h.PrintError("Please specify a project, environment and build ID")
		}
		project := c.Args()[0]
		env := c.Args()[1]
		buildID := c.Args()[2]
		dry := c.Bool("dry")

		response, err := m.RevertEnv(project, env, buildID, cfg.Token, dry)
		if err != nil {
			h.PrintError(err.Error())
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 1 {
			h.PrintError("You should specify the project name")
		}
		if len(c.Args()) < 2 {
			h.PrintError("You should specify the env name")
		}
		project := c.Args()[0]
		env := c.Args()[1]
		if c.String("build") != "" {
			definition, err := m.BuildDefinitionFromIndex(cfg.Token, project, env, c.String("build"))
			if err != nil {
				h.PrintError(err.Error())
			}
			fmt.Println(string(definition))
		} else {
			definition, err := m.LatestBuildDefinition(cfg.Token, project, env)
			if err != nil {
				h.PrintError(err.Error())
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) == 0 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) == 1 {
			h.PrintError("You should specify an existing env name")
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
			h.PrintError(err.Error())
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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 4 {
			h.PrintError("You should specify the project and env names and two build ids to compare them")
		}

		project := c.Args()[0]
		env := c.Args()[1]
		b1 := c.Args()[2]
		b2 := c.Args()[3]

		id1, err := m.BuildIDFromIndex(cfg.Token, project, env, b1)
		if err != nil {
			h.PrintError(err.Error())
		}
		id2, err := m.BuildIDFromIndex(cfg.Token, project, env, b2)
		if err != nil {
			h.PrintError(err.Error())
		}

		def1, err := m.BuildDefinitionByID(cfg.Token, project, env, id1)
		if err != nil {
			h.PrintError(err.Error())
		}

		def2, err := m.BuildDefinitionByID(cfg.Token, project, env, id2)
		if err != nil {
			h.PrintError(err.Error())
		}

		view.PrintEnvDiff(id1, id2, def1, def2)

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
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) == 0 {
			h.PrintError("You should specify an existing project name")
		}
		if len(c.Args()) == 1 {
			h.PrintError("You should specify a valid environment name")
		}

		if c.String("filters") != "" {
			filters = strings.Split(c.String("filters"), ",")
		}

		project := c.Args()[0]
		name := c.Args()[1]
		_, err = m.Import(cfg.Token, name, project, filters)

		if err != nil {
			h.PrintError(err.Error())
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
		SyncEnv,
		ResolveEnv,
		ReviewEnv,
		ScheduleEnv,
	},
}
