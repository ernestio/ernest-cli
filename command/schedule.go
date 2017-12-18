package command

// CmdProject subcommand
import (
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// EnvListSchedules : Gets a list of env schedules
var EnvListSchedules = cli.Command{
	Name:        "list",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.schedules.list.usage"),
	ArgsUsage:   h.T("envs.schedules.list.args"),
	Description: h.T("envs.schedules.list.description"),
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

		list, err := m.EnvSchedules(cfg.Token, project, env)
		if err != nil {
			h.PrintError("Environment does not exist!")
		}

		view.PrintScheduleList(list)
		return nil
	},
}

// EnvAddSchedule : Appends a new schedule to a given environment
var EnvAddSchedule = cli.Command{
	Name:        "add",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.schedules.add.usage"),
	ArgsUsage:   h.T("envs.schedules.add.args"),
	Description: h.T("envs.schedules.add.description"),
	Flags: []cli.Flag{
		stringFlagND("action", "defines what action should be scheduled possible values are [power_on, power_off, sync]"),
		stringFlagND("instance_type", "power_on and power_off accept an instance_type to be powered on an off"),
		stringFlagND("schedule", "sets the automatic schedule. Accepts cron syntax, i.e. '@every 1d', '@weekly' or '0 0 * * * *' (Daily at midnight)"),
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 3 {
			h.PrintError("Required arguments are <project> <env> <schedule_name>")
		}
		project := c.Args()[0]
		env := c.Args()[1]
		name := c.Args()[2]
		for _, flag := range []string{"schedule", "action"} {
			if !c.IsSet(flag) {
				h.PrintError("Missing required flag " + flag)
			}
		}

		e, err := m.EnvStatus(cfg.Token, project, env)
		if err != nil {
			h.PrintError("Environment does not exist!")
		}

		if c.String("action") == "sync" {
			e.Options["sync_interval"] = c.String("schedule")
		} else {
			schedule := make(map[string]string, 0)
			schedule["name"] = name
			schedule["action"] = c.String("action")
			schedule["interval"] = c.String("schedule")
			schedule["instance_type"] = c.String("instance_type")
			if e.Schedules == nil {
				e.Schedules = make(map[string]interface{}, 0)
			}
			e.Schedules[name] = schedule
		}

		err = m.UpdateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c), e.Options, e.Schedules)
		if err != nil {
			h.PrintError(err.Error())
		}

		color.Green("Environment schedules successfully updated")

		return nil
	},
}

// EnvRmSchedule : Appends a new schedule to a given environment
var EnvRmSchedule = cli.Command{
	Name:        "delete",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.schedules.rm.usage"),
	ArgsUsage:   h.T("envs.schedules.rm.args"),
	Description: h.T("envs.schedules.rm.description"),
	Flags:       AllProviderFlags,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) < 3 {
			h.PrintError("Required arguments are <project> <env> <schedule_name>")
		}

		project := c.Args()[0]
		env := c.Args()[1]
		name := c.Args()[2]

		e, err := m.EnvStatus(cfg.Token, project, env)
		if err != nil {
			h.PrintError("Environment does not exist!")
		}

		if name == "sync" {
			e.Options["sync_interval"] = nil
		} else {
			delete(e.Schedules, name)
		}

		err = m.UpdateEnv(cfg.Token, env, project, ProviderFlagsToSlice(c), e.Options, e.Schedules)
		if err != nil {
			h.PrintError(err.Error())
		}

		color.Green("Environment schedules successfully updated")

		return nil
	},
}

// ScheduleEnv ...
var ScheduleEnv = cli.Command{
	Name:    "schedule",
	Aliases: []string{"sch"},
	Usage:   "Scheduling environment related subcommands",
	Subcommands: []cli.Command{
		EnvListSchedules,
		EnvAddSchedule,
		EnvRmSchedule,
	},
}
