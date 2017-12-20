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
		tStringFlag("logger.set.flags.logfile"),
		tStringFlag("logger.set.flags.token"),
		tStringFlag("logger.set.flags.env"),
		tStringFlag("logger.set.flags.hostname"),
		tIntFlag("logger.set.flags.port"),
		tIntFlag("logger.set.flags.timeout"),
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
				h.PrintError(h.T("logger.set.errors.logfile"))
			}
		} else if logger.Type == "logstash" {
			if logger.Hostname == "" {
				h.PrintError(h.T("logger.set.errors.hostname"))
			}
			if logger.Port == 0 {
				h.PrintError(h.T("logger.set.errors.port"))
			}
			if logger.Timeout == 0 {
				h.PrintError(h.T("logger.set.errors.timeout"))
			}
		} else if logger.Type == "rollbar" {
			if logger.Token == "" {
				h.PrintError(h.T("logger.set.errors.token"))
			}
			if logger.Environment == "" {
				logger.Environment = "development"
			}
		} else {
			color.Red(h.T("logger.set.errors.type"))
			return nil
		}

		client.Logger().Create(&logger)
		color.Green(h.T("logger.set.success"))

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
		color.Green(h.T("logger.del.success"))

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
