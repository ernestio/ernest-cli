/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdDatacenter subcommand
import (
	"fmt"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CreateVcloudDatacenter : Creates a VCloud Datacenter
var CreateVcloudDatacenter = cli.Command{
	Name:  "vcloud",
	Usage: "Create a new vcloud datacenter.",
	Description: `Create a new vcloud datacenter on the targeted instance of Ernest.

   Example:
    $ ernest datacenter create vcloud --user username --password xxxx --org MY-ORG-NAME --vse-url http://vse.url --vcloud-url https://myernest.com --public-network MY-PUBLIC-NETWORK mydatacenter

   Template example:
    $ ernest datacenter create vcloud --template mydatacenter.yml mydatacenter
    Where mydatacenter.yaml will look like:
      ---
      fake: true
      org: org
      password: pwd
      public-network: MY-NETWORK
      user: bla
      vcloud-url: "http://ss.com"
      vse-url: "http://ss.com"

	`,
	ArgsUsage: "<datacenter-name>",
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
			Usage: "Datacenter template",
		},
		cli.BoolFlag{
			Name:  "fake",
			Usage: "Fake datacenter",
		},
	},
	Action: func(c *cli.Context) error {
		var errs []string

		if len(c.Args()) == 0 {
			color.Red("You should specify the datacenter name")
			return nil
		}
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		name := c.Args()[0]
		var url, network, user, org, password, username string
		var fake bool

		template := c.String("template")
		if template != "" {
			var t model.DatacenterTemplate
			if err := getDatacenterTemplate(template, &t); err != nil {
				color.Red(err.Error())
				return nil
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
		if fake == false {
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
			color.Red("Please, fix the error shown below to continue")
			for _, e := range errs {
				fmt.Println("  - " + e)
			}
			return nil
		}

		body, err := m.CreateVcloudDatacenter(cfg.Token, name, rtype, username, password, url, network, c.String("vse-url"))
		if err != nil {
			color.Red(body)
		} else {
			color.Green("Datacenter '" + name + "' successfully created ")
		}
		return nil
	},
}

// DeleteDatacenter : Datacenter deletion command definition
var DeleteDatacenter = cli.Command{
	Name:      "delete",
	Usage:     "Deletes the specified datacenter.",
	ArgsUsage: "<datacenter-name>",
	Description: `Deletes the name specified datacenter.

   Example:
    $ ernest datacenter delete my_datacenter
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify the datacenter name")
			return nil
		}
		name := c.Args()[0]

		err := m.DeleteDatacenter(cfg.Token, name)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Datacenter " + name + " successfully removed")

		return nil
	},
}

// UpdateVCloudDatacenter : Updates the specified VCloud datacenter
var UpdateVCloudDatacenter = cli.Command{
	Name:      "vcloud",
	Usage:     "Updates the specified VCloud datacenter.",
	ArgsUsage: "<datacenter-name>",
	Description: `Updates the specified VCloud datacenter.

   Example:
    $ ernest datacenter update vcloud --user <me> --org <org> --password <secret> my_datacenter
	`,
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
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		if len(c.Args()) == 0 {
			color.Red("You should specify the datacenter name")
			return nil
		}
		name := c.Args()[0]
		user = c.String("user")
		password = c.String("password")
		org = c.String("org")

		if user == "" {
			color.Red("You should specify user name with '--user' flag")
			return nil
		}
		if password == "" {
			color.Red("You should specify user password with '--password' flag")
			return nil
		}
		if org == "" {
			color.Red("You should specify user org with '--org' flag")
			return nil
		}

		err := m.UpdateVCloudDatacenter(cfg.Token, name, user+"@"+org, password)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Datacenter " + name + " successfully updated")

		return nil
	},
}
