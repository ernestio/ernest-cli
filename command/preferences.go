/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdDatacenter subcommand
import (
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"

	h "github.com/ernestio/ernest-cli/helper"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// ListLoggers : Lists all loggers confured on ernest
var ListLoggers = cli.Command{
	Name:        "list",
	Usage:       h.T("logger.list.usage"),
	ArgsUsage:   h.T("logger.list.args"),
	Description: h.T("logger.list.description"),
	Action: func(c *cli.Context) error {
		client := esetup(c, NonAdminValidation)

		loggers := client.Logger().List()
		view.PrintLoggerList(loggers)

		return nil
	},
}

// SetLogger : Creates / updates a looger based on it type
var SetLogger = cli.Command{
	Name:        "add",
	Usage:       h.T("logger.set.usage"),
	ArgsUsage:   h.T("logger.set.args"),
	Description: h.T("logger.set.description"),
	Flags: []cli.Flag{
		stringFlagND("logfile", "Specify the path for the loging file"),
		stringFlagND("token", "Rollbar token"),
		stringFlagND("env", "Rollbar environment"),
		stringFlagND("hostname", "Logstash hostname"),
		intFlag("port", "Logstash port"),
		intFlag("timeout", "Logstash timeout"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "logger.set.args")
		client := esetup(c, NonAdminValidation)

		logger := emodels.Logger{
			Type:        c.Args()[0],
			Logfile:     c.String("logfile"),
			Hostname:    c.String("hostname"),
			Port:        c.Int("port"),
			Timeout:     c.Int("timeout"),
			Token:       c.String("token"),
			Environment: c.String("env"),
		}
		if logger.Type == "basic" {
			if logger.Logfile == "" {
				h.PrintError("You should specify a logfile with --logfile flag")
			}
		} else if logger.Type == "logstash" {
			if logger.Hostname == "" {
				h.PrintError("You should specify a logstash hostname  with --hostname flag")
			}
			if logger.Port == 0 {
				h.PrintError("You should specify a logstash port with --port flag")
			}
			if logger.Timeout == 0 {
				h.PrintError("You should specify a logstash timeout with --timeout flag")
			}
		} else if logger.Type == "rollbar" {
			if logger.Token == "" {
				h.PrintError("You should specify a rollbar token with --token flag")
			}
			if logger.Environment == "" {
				logger.Environment = "development"
			}
		} else {
			color.Red("Invalid type, valid types are basic and logstash")
			return nil
		}

		client.Logger().Create(&logger)
		color.Green("Logger successfully set up")

		return nil
	},
}

// DelLogger : deletes a looger based on it type
var DelLogger = cli.Command{
	Name:        "delete",
	Usage:       h.T("logger.del.usage"),
	ArgsUsage:   h.T("logger.del.args"),
	Description: h.T("logger.del.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "logger.del.args")
		client := esetup(c, NonAdminValidation)

		client.Logger().Delete(c.Args()[0])
		color.Green("Logger successfully deleted")

		return nil
	},
}

// LoggerCommands : Logger related commands
var LoggerCommands = cli.Command{
	Name:        "logger",
	Usage:       "Setup logger preferences.",
	Description: "Setup ernest logger preferenres.",
	Subcommands: []cli.Command{
		ListLoggers,
		SetLogger,
		DelLogger,
	},
}

// CmdPreferences : Preferences setup
var CmdPreferences = cli.Command{
	Name:  "preferences",
	Usage: "Ernest preferences",
	Subcommands: []cli.Command{
		LoggerCommands,
	},
}
