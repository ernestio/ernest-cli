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

// CreateAWSDatacenter : Creates an AWS datacenter
var CreateAWSDatacenter = cli.Command{
	Name:  "aws",
	Usage: "Create a new aws datacenter.",
	Description: `Create a new AWS datacenter on the targeted instance of Ernest.

	Example:
	 $ ernest datacenter create aws --region us-west-2 --access_key_id AKIAIOSFODNN7EXAMPLE --secret_access_key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY my_datacenter

   Template example:
    $ ernest datacenter create aws --template mydatacenter.yml mydatacenter
    Where mydatacenter.yaml will look like:
      ---
      fake: true
      access_key_id : AKIAIOSFODNN7EXAMPLE
      secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
      region: us-west-2
	 `,
	ArgsUsage: "<datacenter-name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "region, r",
			Value: "",
			Usage: "Datacenter region",
		},
		cli.StringFlag{
			Name:  "access_key_id, k",
			Value: "",
			Usage: "AWS access key id",
		},
		cli.StringFlag{
			Name:  "secret_access_key, s",
			Value: "",
			Usage: "AWS Secret access key",
		},
		cli.StringFlag{
			Name:  "template, t",
			Value: "",
			Usage: "Datacenter template",
		},
		cli.BoolFlag{
			Name:  "fake, f",
			Usage: "Fake datacenter",
		},
	},
	Action: func(c *cli.Context) error {
		var errs []string
		var accessKeyID, secretAccessKey, region string
		var fake bool
		m, cfg := setup(c)

		if len(c.Args()) < 1 {
			msg := "You should specify the datacenter name"
			color.Red(msg)
			return nil
		}

		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		name := c.Args()[0]

		template := c.String("template")
		if template != "" {
			var t model.DatacenterTemplate
			if err := getDatacenterTemplate(template, &t); err != nil {
				color.Red(err.Error())
				return nil
			}
			accessKeyID = t.Token
			secretAccessKey = t.Secret
			region = t.Region
			fake = t.Fake
		}
		if c.String("secret_access_key") != "" {
			secretAccessKey = c.String("secret_access_key")
		}
		if c.String("access_key_id") != "" {
			accessKeyID = c.String("access_key_id")
		}
		if c.String("region") != "" {
			region = c.String("region")
		}
		if fake == false {
			fake = c.Bool("fake")
		}

		if secretAccessKey == "" {
			errs = append(errs, "Specify a valid secret access key with --secret_access_key flag")
		}

		if accessKeyID == "" {
			errs = append(errs, "Specify a valid access key id with --access_key_id flag")
		}

		if region == "" {
			errs = append(errs, "Specify a valid region with --region flag")
		}

		if len(errs) > 0 {
			color.Red("Please, fix the error shown below to continue")
			for _, e := range errs {
				fmt.Println("  - " + e)
			}
			return nil
		}

		rtype := "aws"

		if fake {
			rtype = "aws-fake"
		}
		body, err := m.CreateAWSDatacenter(cfg.Token, name, rtype, region, accessKeyID, secretAccessKey)
		if err != nil {
			color.Red(body)
		} else {
			color.Green("Datacenter '" + name + "' successfully created ")
		}
		return nil
	},
}

// UpdateAWSDatacenter : Updates the specified VCloud datacenter
var UpdateAWSDatacenter = cli.Command{
	Name:      "aws",
	Usage:     "Updates the specified AWS datacenter.",
	ArgsUsage: "<datacenter-name>",
	Description: `Updates the specified AWS datacenter.

   Example:
		$ ernest datacenter update aws --access_key_id AKIAIOSFODNN7EXAMPLE --secret_access_key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY my_datacenter
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "access_key_id",
			Value: "",
			Usage: "Your AWS access key id",
		},
		cli.StringFlag{
			Name:  "secret_access_key",
			Value: "",
			Usage: "Your AWS secret access key",
		},
	},
	Action: func(c *cli.Context) error {
		var accessKeyID, secretAccessKey string
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
		accessKeyID = c.String("access_key_id")
		secretAccessKey = c.String("secret_access_key")

		if accessKeyID == "" {
			color.Red("You should specify your aws access key id with '--access_key_id' flag")
			return nil
		}
		if secretAccessKey == "" {
			color.Red("You should specify your aws secret access key with '--secret_access_key' flag")
			return nil
		}

		err := m.UpdateAWSDatacenter(cfg.Token, name, accessKeyID, secretAccessKey)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Datacenter " + name + " successfully updated")

		return nil
	},
}
