/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"fmt"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CreateVcloudProject : Creates a VCloud Project
var CreateVcloudProject = cli.Command{
	Name:        "vcloud",
	Usage:       h.T("vcloud.create.usage"),
	ArgsUsage:   h.T("vcloud.create.args"),
	Description: h.T("vcloud.create.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Your VCloud valid user name",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Your VCloud valid password",
		},
		cli.StringFlag{
			Name:  "org",
			Value: "",
			Usage: "Your vCloud Organization",
		},
		cli.StringFlag{
			Name:  "vse-url",
			Value: "",
			Usage: "VSE URL",
		},
		cli.StringFlag{
			Name:  "vcloud-url",
			Value: "",
			Usage: "VCloud URL",
		},
		cli.StringFlag{
			Name:  "public-network",
			Value: "",
			Usage: "Public Network",
		},
		cli.StringFlag{
			Name:  "template",
			Value: "",
			Usage: "Project template",
		},
		cli.BoolFlag{
			Name:  "fake",
			Usage: "Fake project",
		},
	},
	Action: func(c *cli.Context) error {
		var errs []string

		if len(c.Args()) == 0 {
			h.PrintError("You should specify the project name")
		}
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		name := c.Args()[0]
		var url, network, user, org, password, username string
		var fake bool

		template := c.String("template")
		if template != "" {
			var t model.ProjectTemplate
			if err := getProjectTemplate(template, &t); err != nil {
				h.PrintError(err.Error())
			}
			url = t.URL
			network = t.Network
			user = t.User
			org = t.Org
			password = t.Password
			fake = t.Fake
		}
		if c.String("vcloud-url") != "" {
			url = c.String("vcloud-url")
		}
		if c.String("public-network") != "" {
			network = c.String("public-network")
		}
		if c.String("user") != "" {
			user = c.String("user")
		}
		if c.String("org") != "" {
			org = c.String("org")
		}
		if c.String("password") != "" {
			password = c.String("password")
		}
		if !fake {
			fake = c.Bool("fake")
		}
		username = user + "@" + org

		if url == "" {
			errs = append(errs, "Specify a valid VCloud URL with --vcloud-url flag")
		}
		if network == "" {
			errs = append(errs, "Specify a valid public network with --public-network flag")
		}
		if user == "" {
			errs = append(errs, "Specify a valid user name with --user")
		}
		if org == "" {
			errs = append(errs, "Specify a valid organization with --org")
		}
		if password == "" {
			errs = append(errs, "Specify a valid password with --password")
		}
		rtype := "vcloud"
		if fake {
			rtype = "vcloud-fake"
		}
		if len(errs) > 0 {
			h.PrintError("Please, fix the error shown below to continue")
			for _, e := range errs {
				fmt.Println("  - " + e)
			}
			return nil
		}

		body, err := m.CreateVcloudProject(cfg.Token, name, rtype, username, password, url, network, c.String("vse-url"))
		if err != nil {
			h.PrintError(body)
		} else {
			color.Green("Project '" + name + "' successfully created ")
		}
		return nil
	},
}

// DeleteProject : Project deletion command definition
var DeleteProject = cli.Command{
	Name:        "delete",
	Usage:       h.T("vcloud.delete.usage"),
	ArgsUsage:   h.T("vcloud.delete.args"),
	Description: h.T("vcloud.delete.description"),
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}

		if len(c.Args()) == 0 {
			h.PrintError("You should specify the project name")
		}
		name := c.Args()[0]

		err := m.DeleteProject(cfg.Token, name)
		if err != nil {
			h.PrintError(err.Error())
		}
		color.Green("Project " + name + " successfully removed")

		return nil
	},
}

// UpdateVCloudProject : Updates the specified VCloud project
var UpdateVCloudProject = cli.Command{
	Name:        "vcloud",
	Usage:       h.T("vcloud.update.usage"),
	ArgsUsage:   h.T("vcloud.update.args"),
	Description: h.T("vcloud.update.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Your VCloud valid user name",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Your VCloud valid password",
		},
		cli.StringFlag{
			Name:  "org",
			Value: "",
			Usage: "Your vCloud Organization",
		},
	},
	Action: func(c *cli.Context) error {
		var user, password, org string
		m, cfg := setup(c)
		if cfg.Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}
		if len(c.Args()) == 0 {
			h.PrintError("You should specify the project name")
		}
		name := c.Args()[0]
		user = c.String("user")
		password = c.String("password")
		org = c.String("org")

		if user == "" {
			h.PrintError("You should specify user name with '--user' flag")
		}
		if password == "" {
			h.PrintError("You should specify user password with '--password' flag")
		}
		if org == "" {
			h.PrintError("You should specify user org with '--org' flag")
		}

		err := m.UpdateVCloudProject(cfg.Token, name, user+"@"+org, password)
		if err != nil {
			h.PrintError(err.Error())
		}
		color.Green("Project " + name + " successfully updated")

		return nil
	},
}
