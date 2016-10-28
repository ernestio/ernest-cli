/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdDatacenter subcommand
import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListLoggers : Lists all loggers confured on ernest
var ListLoggers = cli.Command{
	Name:      "list",
	Usage:     "Lists active loggers.",
	ArgsUsage: " ",
	Description: `List active loggers.

   Example:
    $ ernest preferences logger list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		loggers, err := m.ListLoggers(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

		printLoggerList(loggers)

		return nil
	},
}

// SetLogger : Creates / updates a looger based on it type
var SetLogger = cli.Command{
	Name:      "add",
	Usage:     "Creates / updates a logger based on its type.",
	ArgsUsage: " ",
	Description: `Creates / updates a logger based on its types.

   Example:
    $ ernest preferences logger add basic --logfile /tmp/ernest.log
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "logfile",
			Usage: "Specify the path for the loging file",
		},
		cli.StringFlag{
			Name:  "hostname",
			Usage: "Logstash hostname",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "Logstash port",
		},
		cli.IntFlag{
			Name:  "timeout",
			Usage: "Logstash timeout",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify the logger type (basic | logstash)")
			return nil
		}

		logger := Logger{
			Type:     c.Args()[0],
			Logfile:  c.String("logfile"),
			Hostname: c.String("hostname"),
			Port:     c.Int("port"),
			Timeout:  c.Int("timeout"),
		}
		if logger.Type == "basic" {
			if logger.Logfile == "" {
				color.Red("You should specify a logfile with --logfile flag")
				return nil
			}
		} else if logger.Type == "generic" {
			if logger.Hostname == "" {
				color.Red("You should specify a logstash hostname  with --hostname flag")
				return nil
			}
			if logger.Port == 0 {
				color.Red("You should specify a logstash port with --port flag")
				return nil
			}
			if logger.Timeout == 0 {
				color.Red("You should specify a logstash timeout with --timeout flag")
				return nil
			}
		} else {
			color.Red("Invalid type, valid types are basic and logstash")
			return nil
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		err := m.SetLogger(cfg.Token, logger)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

		color.Green("Logger successfully set up")

		return nil
	},
}

// DelLogger : deletes a looger based on it type
var DelLogger = cli.Command{
	Name:      "delete",
	Usage:     "Deletes a logger based on its type.",
	ArgsUsage: " ",
	Description: `Deletes a logger based on its types.

   Example:
    $ ernest preferences logger delete basic
	`,
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify the logger type (basic | logstash)")
			return nil
		}

		logger := Logger{
			Type: c.Args()[0],
		}

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		err := m.DelLogger(cfg.Token, logger)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

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
