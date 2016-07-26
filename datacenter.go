/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdDatacenter subcommand
import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListDatacenters ...
var ListDatacenters = cli.Command{
	Name:      "list",
	Usage:     "List available datacenters.",
	ArgsUsage: " ",
	Description: `List available datacenters.

   Example:
    $ ernest datacenter list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		datacenters, err := m.ListDatacenters(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		fmt.Fprintln(w, "NAME\tTYPE")
		for _, datacenter := range datacenters {
			str := fmt.Sprintf("%s\t%s", datacenter.Name, datacenter.Type)
			fmt.Fprintln(w, str)
		}
		w.Flush()
		return nil
	},
}

// CreateAWSDatacenter ...
var CreateAWSDatacenter = cli.Command{
	Name:  "aws",
	Usage: "Create a new aws datacenter.",
	Description: `Create a new AWS datacenter on the targeted instance of Ernest.

	Example:
	 $ ernest datacenter create aws --region region --token token --secret secret my_datacenter
	 `,
	ArgsUsage: "<datacenter-name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "region",
			Value: "",
			Usage: "Datacenter region",
		},
		cli.StringFlag{
			Name:  "token",
			Value: "",
			Usage: "AWS Token",
		},
		cli.StringFlag{
			Name:  "secret",
			Value: "",
			Usage: "AWS Secret",
		},
		cli.BoolFlag{
			Name:  "fake",
			Usage: "Fake datacenter",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify the datacenter name"
			color.Red(msg)
			return errors.New(msg)
		}

		name := c.Args()[0]

		token := c.String("token")
		if token == "" {
			msg := "Token not specified"
			color.Red(msg)
			return errors.New(msg)
		}

		secret := c.String("secret")
		if secret == "" {
			msg := "Secret not specified"
			color.Red(msg)
			return errors.New(msg)
		}

		region := c.String("region")
		if region == "" {
			msg := "Region not specified"
			color.Red(msg)
			return errors.New(msg)
		}

		rtype := "aws"

		if c.Bool("fake") {
			rtype = "aws-fake"
		}
		m, cfg := setup(c)
		_, err := m.CreateAWSDatacenter(cfg.Token, name, rtype, region, token, secret)
		if err != nil {
			color.Red(err.Error())
		}
		return nil
	},
}

// CreateVcloudDatacenter ...
var CreateVcloudDatacenter = cli.Command{
	Name:  "vcloud",
	Usage: "Create a new vcloud datacenter.",
	Description: `Create a new vcloud datacenter on the targeted instance of Ernest.

   Example:
    $ ernest datacenter create vcloud --datacenter-user username --datacenter-password xxxx --datacenter-org MY-ORG-NAME --vse-url http://vse.url mydatacenter https://myernest.com MY-PUBLIC-NETWORK
	`,
	ArgsUsage: "<datacenter-name> <vcloud-url> <network-name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "datacenter-user",
			Value: "",
			Usage: "User to be configured with the datacenter",
		},
		cli.StringFlag{
			Name:  "datacenter-password",
			Value: "",
			Usage: "Password related with user",
		},
		cli.StringFlag{
			Name:  "datacenter-org",
			Value: "",
			Usage: "vCloud Organization name",
		},
		cli.StringFlag{
			Name:  "vse-url",
			Value: "",
			Usage: "VSE URL",
		},
		cli.BoolFlag{
			Name:  "fake",
			Usage: "Fake datacenter",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 3 {
			msg := "You should specify the datacenter name, vcloud URL and network name"
			color.Red(msg)
			return errors.New(msg)
		}
		name := c.Args()[0]
		user := c.String("datacenter-user") + "@" + c.String("datacenter-org")
		if user == "" {
			msg := "User not specified"
			color.Red(msg)
			return errors.New(msg)
		}
		password := c.String("datacenter-password")
		if password == "" {
			msg := "Password not specified"
			color.Red(msg)
			return errors.New("Password not specified")
		}
		rtype := "vcloud"
		if c.Bool("fake") {
			rtype = "vcloud-fake"
		}
		m, cfg := setup(c)
		_, err := m.CreateVcloudDatacenter(cfg.Token, name, rtype, user, password, c.Args()[1], c.Args()[2], c.String("vse-url"))
		if err != nil {
			color.Red(err.Error())
		}
		return nil
	},
}

// CreateDatacenters ...
var CreateDatacenters = cli.Command{
	Name:        "create",
	Usage:       "Create a new datacenter.",
	Description: "Create a new datacenter on the targeted instance of Ernest.",
	Subcommands: []cli.Command{
		CreateVcloudDatacenter,
		CreateAWSDatacenter,
	},
}

// CmdDatacenter ...
var CmdDatacenter = cli.Command{
	Name:  "datacenter",
	Usage: "Datacenter related subcommands",
	Subcommands: []cli.Command{
		ListDatacenters,
		CreateDatacenters,
	},
}
