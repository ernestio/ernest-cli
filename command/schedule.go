package command

// CmdProject subcommand
import (
	"strings"

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
		tStringFlagND("envs.schedule.add.flags.resolution"),
		tStringFlagND("envs.schedule.add.flags.instances"),
		tStringFlagND("envs.schedule.add.flags.schedule"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 3, "envs.schedule.list.args")
		client := esetup(c, AuthUsersValidation)

		env := client.Environment().Get(c.Args()[0], c.Args()[1])

		schedule := make(map[string]interface{}, 0)
		schedule["name"] = c.Args()[2]
		schedule["type"] = c.String("action")
		schedule["interval"] = c.String("schedule")

		switch c.String("action") {
		case "sync":
			resolution := c.String("resolution")
			if resolution == "" {
				resolution = "manual"
			}
			schedule["resolution"] = resolution
		case "power_on", "power_off":
			schedule["instances"] = strings.Split(c.String("instances"), ",")
		default:
			h.PrintError("unsupported action type: " + c.String("action"))
		}

		if env.Schedules == nil {
			env.Schedules = make(map[string]interface{}, 0)
		}

		env.Schedules[c.Args()[2]] = schedule
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
		delete(env.Schedules, c.Args()[2])

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
