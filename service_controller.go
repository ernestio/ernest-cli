/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdDatacenter subcommand
import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListServices ...
var ListServices = cli.Command{
	Name:      "list",
	Usage:     "List available services.",
	ArgsUsage: " ",
	Description: `List available services and shows its most relevant information.

   Example:
    $ ernest service list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		services, err := m.ListServices(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		printServiceList(services)
		return nil
	},
}

// ApplyService command
// Applies changes described on a YAML file to a service
var ApplyService = cli.Command{
	Name:      "apply",
	Aliases:   []string{"a"},
	ArgsUsage: "<file.yml>",
	Usage:     "Builds or changes infrastructure.",
	Description: `Sends a service YAML description file to Ernest to be executed.
   You must be logged in to execute this command.

   If the file is not provided, ernest.yml will be used by default.

   Example:
    $ ernest apply myservice.yml
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
		_, err := m.Apply(cfg.Token, file, true)
		if err != nil {
			color.Red(err.Error())
		}
		return nil
	},
}

// DestroyService command
var DestroyService = cli.Command{
	Name:      "destroy",
	Aliases:   []string{"d"},
	ArgsUsage: "<service_name>",
	Usage:     "Destroy a service.",
	Description: `Destroys a service by its name.

   Example:
    $ ernest destroy myservice
  `,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force,f",
			Usage: "Force destroy command without asking for permission.",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify an existing service name")
			return nil
		}
		name := c.Args()[0]

		if c.Bool("force") {
			err := m.Destroy(cfg.Token, name, true)
			if err != nil {
				color.Red(err.Error())
			}
		} else {
			fmt.Print("Are you sure? Please type yes or no and then press enter: ")
			if askForConfirmation() {
				err := m.Destroy(cfg.Token, name, true)
				if err != nil {
					color.Red(err.Error())
					return err
				}
			}
		}
		return nil
	},
}

// HistoryService command
// Shows the history of a service, a list of builds
var HistoryService = cli.Command{
	Name:      "history",
	Usage:     "Shows the history of a service, a list of builds",
	ArgsUsage: "<service_name>",
	Description: `Shows the history of a service, a list of builds and its status and basic information.

   Example:
    $ ernest history myservice
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify an existing service name")
			return nil
		}

		serviceName := c.Args()[0]

		services, _ := m.ListBuilds(serviceName, cfg.Token)
		printServiceHistory(services)
		return nil
	},
}

// ResetService command
var ResetService = cli.Command{
	Name:      "reset",
	ArgsUsage: "<service_name>",
	Usage:     "Reset an in progress service.",
	Description: `Reseting a service creation may cause problems, please make sure you know what are you doing.

   Example:
    $ ernest reset myservice
  `,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify the service name")
		}
		serviceName := c.Args()[0]
		err := m.ResetService(serviceName, cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Red("You've successfully resetted the service '" + serviceName + "'")

		return nil
	},
}

// DefinitionService command
// Shows the current definition of a service by its name
var DefinitionService = cli.Command{
	Name:      "definition",
	Aliases:   []string{"s"},
	ArgsUsage: "<service_name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "build",
			Value: "",
			Usage: "Build ID",
		},
	},
	Usage: "Show the current definition of a service by its name",
	Description: `Show the current definition of a service by its name getting the definition about the build.

   Example:
    $ ernest service definition myservice
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) < 1 {
			color.Red("You should specify the service name")
			return nil
		}
		serviceName := c.Args()[0]
		if c.String("build") != "" {
			service, err := m.ServiceBuildStatus(cfg.Token, serviceName, c.String("build"))
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}
			fmt.Println(service.Definition)
		} else {
			service, err := m.ServiceStatus(cfg.Token, serviceName)
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}

			fmt.Println(service.Definition)
		}
		return nil
	},
}

// InfoService : Shows detailed information of a service by its name
var InfoService = cli.Command{
	Name:      "info",
	Aliases:   []string{"i"},
	ArgsUsage: "<service_name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "build",
			Value: "",
			Usage: "Build ID",
		},
	},
	Usage: "$ ernest service info <my_service> --build <specific build>",
	Description: `Will show detailed information of the last build of a specified service.
	In case you specify --build option you will be able to output the detailed information of specific build of a service.

   Examples:
    $ ernest service definition myservice
    $ ernest service definition myservice --build build1
	`,
	Action: func(c *cli.Context) error {
		var err error
		var service Service

		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify an existing service name")
			return nil
		}

		name := c.Args()[0]
		if c.String("build") != "" {
			build := c.String("build")
			service, err = m.ServiceBuildStatus(cfg.Token, name, build)
		} else {
			service, err = m.ServiceStatus(cfg.Token, name)
		}

		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		printServiceInfo(&service)
		return nil
	},
}

// CmdService ...
var CmdService = cli.Command{
	Name:  "service",
	Usage: "Service related subcommands",
	Subcommands: []cli.Command{
		ListServices,
		ApplyService,
		DestroyService,
		HistoryService,
		ResetService,
		DefinitionService,
		InfoService,
		MonitorService,
	},
}
