/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdProject subcommand
import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CreateAzureProject : Creates an AWS project
var CreateAzureProject = cli.Command{
	Name:  "azure",
	Usage: "Create a new azure project.",
	Description: `Create a new Azure project on the targeted instance of Ernest.

	Example:
	 $ ernest project create azure --region westus --subscription_id SUBSCRIPTION --client_id USER --client_secret PASSWORD --tenant_id TENANT --environment public my_project

   Template example:
    $ ernest project create azure --template myproject.yml myproject
    Where myproject.yaml will look like:
      ---
      fake: true
      region: westus
			subscription_id: SUBSCRIPTION
			client_id: USER
			client_secret: PASSWORD
			tenant_id: TENANT
			environment: public
	 `,
	ArgsUsage: "<project-name>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "region, r",
			Value: "",
			Usage: "Project region",
		},
		cli.StringFlag{
			Name:  "subscription_id, s",
			Value: "",
			Usage: "Azure subscription id",
		},
		cli.StringFlag{
			Name:  "client_id, c",
			Value: "",
			Usage: "Azure client id",
		},
		cli.StringFlag{
			Name:  "client_secret, p",
			Value: "",
			Usage: "Azure client secret",
		},
		cli.StringFlag{
			Name:  "tenant_id, t",
			Value: "",
			Usage: "Azure tenant_id",
		},
		cli.StringFlag{
			Name:  "environment, e",
			Value: "",
			Usage: "Azure environment. Supported values are public(default), usgovernment, german and chine",
		},
		cli.BoolFlag{
			Name:  "fake, f",
			Usage: "Fake project",
		},
	},
	Action: func(c *cli.Context) error {
		var errs []string
		var region, subscriptionID, clientID, clientSecret, tenantID, environment string
		var fake bool
		m, cfg := setup(c)

		if len(c.Args()) < 1 {
			msg := "You should specify the project name"
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
			/*
				var t model.ProjectTemplate
				if err := getProjectTemplate(template, &t); err != nil {
					color.Red(err.Error())
					return nil
				}
				accessKeyID = t.Token
				secretAccessKey = t.Secret
				region = t.Region
				fake = t.Fake
			*/
		}
		if c.String("region") != "" {
			region = c.String("region")
		}
		if c.String("subscription_id") != "" {
			subscriptionID = c.String("subscription_id")
		}
		if c.String("client_id") != "" {
			clientID = c.String("client_id")
		}
		if c.String("client_secret") != "" {
			clientSecret = c.String("client_secret")
		}
		if c.String("tenant_id") != "" {
			tenantID = c.String("tenant_id")
		}
		if c.String("environment") != "" {
			environment = c.String("environment")
		}
		if !fake {
			fake = c.Bool("fake")
		}

		if subscriptionID == "" {
			errs = append(errs, "Specify a valid subscription id with --subscription_id flag")
		}
		if clientID == "" {
			errs = append(errs, "Specify a valid client id with --client_id flag")
		}
		if region == "" {
			errs = append(errs, "Specify a valid region with --region flag")
		}
		if tenantID == "" {
			errs = append(errs, "Specify a valid tenant id with --tenant_id flag")
		}
		if environment == "" {
			environment = "public"
		}

		if len(errs) > 0 {
			color.Red("Please, fix the error shown below to continue")
			for _, e := range errs {
				fmt.Println("  - " + e)
			}
			return nil
		}

		rtype := "azure"

		if fake {
			rtype = "azure-fake"
		}
		body, err := m.CreateAzureProject(cfg.Token, name, rtype, region, subscriptionID, clientID, clientSecret, tenantID, environment)
		if err != nil {
			color.Red(body)
		} else {
			color.Green("Project '" + name + "' successfully created ")
		}
		return nil
	},
}

// UpdateAzureProject : Updates the specified VCloud project
var UpdateAzureProject = cli.Command{
	Name:      "azure",
	Usage:     "Updates the specified Azure project.",
	ArgsUsage: "<project-name>",
	Description: `Updates the specified Azure project.

   Example:
	 $ ernest project update azure --subscription_id SUBSCRIPTION --client_id USER --client_secret PASSWORD --tenant_id TENANT --environment public my_project
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "subscription_id, s",
			Value: "",
			Usage: "Azure subscription id",
		},
		cli.StringFlag{
			Name:  "client_id, c",
			Value: "",
			Usage: "Azure client id",
		},
		cli.StringFlag{
			Name:  "client_secret, p",
			Value: "",
			Usage: "Azure client secret",
		},
		cli.StringFlag{
			Name:  "tenant_id, t",
			Value: "",
			Usage: "Azure tenant_id",
		},
		cli.StringFlag{
			Name:  "environment, e",
			Value: "",
			Usage: "Azure environment. Supported values are public(default), usgovernment, german and chine",
		},
	},
	Action: func(c *cli.Context) error {
		var errs []string
		var subscriptionID, clientID, clientSecret, tenantID, environment string
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify the project name")
			return nil
		}
		name := c.Args()[0]
		if c.String("subscription_id") != "" {
			subscriptionID = c.String("subscription_id")
		}
		if c.String("client_id") != "" {
			clientID = c.String("client_id")
		}
		if c.String("client_secret") != "" {
			clientSecret = c.String("client_secret")
		}
		if c.String("tenant_id") != "" {
			tenantID = c.String("tenant_id")
		}
		if c.String("environment") != "" {
			environment = c.String("environment")
		}

		if subscriptionID == "" {
			errs = append(errs, "Specify a valid subscription id with --subscription_id flag")
		}
		if clientID == "" {
			errs = append(errs, "Specify a valid client id with --client_id flag")
		}
		if clientSecret == "" {
			errs = append(errs, "Specify a valid client secret with --client_secret flag")
		}
		if tenantID == "" {
			errs = append(errs, "Specify a valid tenant id with --tenant_id flag")
		}
		if environment == "" {
			errs = append(errs, "Specify a valid environment with --environment flag")
		}

		if len(errs) > 0 {
			color.Red("Please, fix the error shown below to continue")
			for _, e := range errs {
				fmt.Println("  - " + e)
			}
			return nil
		}

		err := m.UpdateAzureProject(cfg.Token, name, subscriptionID, clientID, clientSecret, tenantID, environment)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Project " + name + " successfully updated")

		return nil
	},
}
