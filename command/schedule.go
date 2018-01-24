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
	Usage:       h.T("envs.schedule.list.usage"),
	ArgsUsage:   h.T("envs.schedule.list.args"),
	Description: h.T("envs.schedule.list.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "envs.schedule.list.args")
		client := esetup(c, AuthUsersValidation)

		env := client.Environment().Get(c.Args()[0], c.Args()[1])
		list := env.Schedules

		view.PrintScheduleList(list)
		return nil
	},
}

// EnvAddSchedule : Appends a new schedule to a given environment
var EnvAddSchedule = cli.Command{
	Name:        "add",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.schedule.add.usage"),
	ArgsUsage:   h.T("envs.schedule.add.args"),
	Description: h.T("envs.schedule.add.description"),
	Flags: []cli.Flag{
		tStringFlagND("envs.schedule.add.flags.action"),
		tStringFlagND("envs.schedule.add.flags.instance_type"),
		tStringFlagND("envs.schedule.add.flags.schedule"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "envs.schedule.list.args")
		client := esetup(c, AuthUsersValidation)

		env := client.Environment().Get(c.Args()[0], c.Args()[1])

		if c.String("action") == "sync" {
			env.Options["sync_interval"] = c.String("schedule")
		} else {
			schedule := make(map[string]string, 0)
			schedule["name"] = c.Args()[2]
			schedule["action"] = c.String("action")
			schedule["interval"] = c.String("schedule")
			schedule["instance_type"] = c.String("instance_type")
			if env.Schedules == nil {
				env.Schedules = make(map[string]interface{}, 0)
			}
			env.Schedules[c.Args()[2]] = schedule
		}
		client.Environment().Update(env)
		color.Green(h.T("envs.schedule.add.success"))

		return nil
	},
}

// EnvRmSchedule : Appends a new schedule to a given environment
var EnvRmSchedule = cli.Command{
	Name:        "delete",
	Aliases:     []string{"a"},
	Usage:       h.T("envs.schedule.rm.usage"),
	ArgsUsage:   h.T("envs.schedule.rm.args"),
	Description: h.T("envs.schedule.rm.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "envs.schedule.list.args")
		client := esetup(c, AuthUsersValidation)

		env := client.Environment().Get(c.Args()[0], c.Args()[1])
		if c.Args()[2] == "sync" {
			env.Options["sync_interval"] = nil
		} else {
			delete(env.Schedules, c.Args()[2])
		}
		client.Environment().Update(env)
		color.Green(h.T("envs.schedule.rm.success"))

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
