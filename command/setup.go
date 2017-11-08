/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"fmt"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/manager"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

// CmdSetup : Setup an ernest instance
var CmdSetup = cli.Command{
	Name:        "setup",
	Usage:       h.T("setup.usage"),
	Description: h.T("setup.description"),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: "",
			Usage: "Admin user",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Admin password",
		},
		cli.StringFlag{
			Name:  "target, t",
			Value: "",
			Usage: "Ernest location",
		},
	},
	Action: func(c *cli.Context) error {
		var token string
		var err error

		color.Blue("This process will guide you to setting up your ernest instance")
		m, cfg := setup(c)
		cfg.URL = c.String("target")
		if cfg.URL == "" {
			cfg.URL = reAskDef("- Introduce the location of your ernest instance (https://ernest.local/):", "https://ernest.local/")
		}
		_ = persistTarget(cfg)
		color.Green("Target set to '" + cfg.URL + "'")

		adminToken := adminLogin(m, c.String("user"), c.String("password"))

		usr, pwd := createUser(adminToken, m)

		// Login as plain user
		if token, err = m.Login(usr, pwd, ""); err != nil {
			fmt.Println("Ups, something went wrong creating the plain user")
			return nil
		}

		cfg.Token = token
		cfg.User = usr
		if err := model.SaveConfig(cfg); err == nil {
			color.Green("You're now logged in as '" + usr + "'")
		}
		color.Green("Ernest is successfully set up!")
		fmt.Println("")
		fmt.Print("Would you like to create a project now? (Y/n): ")

		projectCreation := askForConfirmation()
		if projectCreation {
			createProject(token, m)
		}
		color.Green("Congratulations, you've successfully set up your ernest istance.")
		fmt.Println("")
		color.Green("What's next?")

		if !projectCreation {
			color.Green("- Create a project with command 'ernest-cli project create --help'")
		}
		color.Green("- Apply an env definition with 'ernest-cli env apply your-template.yml'")
		color.Green("- Need some extra documentation? Run 'ernest-cli docs'")

		return nil
	},
}

func createProject(token string, m *manager.Manager) {
	dt := reAsk("- Please specify the provider you want to work with? (vcloud / aws)")
	if dt == "aws" {
		createAWSProject(token, m)
	} else if dt == "vcloud" {
		createVCloudProject(token, m)
	} else {
		fmt.Println("Invalid provider. Valid providers are (vcloud and aws)")
		createProject(token, m)
	}
}

func createAWSProject(token string, m *manager.Manager) {
	fmt.Println("In order to create an AWS project we will need some info: ")
	dt := "aws"
	dname := reAsk("- Name:")
	dregion := reAsk("- Region:")
	dkey := reAsk("- Access key id:")
	dsecret := reAsk("- Secret access key:")
	if body, err := m.CreateAWSProject(token, dname, dt, dregion, dkey, dsecret); err != nil {
		color.Red("ERROR: " + body)
		createAWSProject(token, m)
	}
}

func createVCloudProject(token string, m *manager.Manager) {
	fmt.Println("In order to create a VCloud project we will need some info: ")
	dt := "vcloud"
	dname := reAsk("- Name:")
	dusername := reAsk("- Username:")
	fmt.Print("- Password: ")
	pass, _ := gopass.GetPasswdMasked()
	dpassword := string(pass)
	durl := reAsk("- URL:")
	dnetwork := reAsk("- Network:")
	dvse := reAsk("- VSE url:")

	if body, err := m.CreateVcloudProject(token, dname, dt, dusername, dpassword, durl, dnetwork, dvse); err != nil {
		color.Red("ERROR: " + body)
		createVCloudProject(token, m)
	}
}

func createUser(adminToken string, m *manager.Manager) (string, string) {
	var err error

	fmt.Println("")
	fmt.Println("In order to perform basic actions, ernest needs a basic account.")
	fmt.Println("Please introduce your account information:")
	usr := reAsk("- New username:")
	pwd := getConfirmedPasswords()
	// TODO : Check success?
	if err = m.CreateUser(adminToken, usr, "", usr, pwd); err != nil {
		color.Red(err.Error() + ". Please try again")
		return createUser(adminToken, m)
	}
	color.Green("User '" + usr + "' successfully created!")

	return usr, pwd
}

func getConfirmedPasswords() string {
	fmt.Print("- Password: ")
	pass, _ := gopass.GetPasswdMasked()
	pwd := string(pass)
	fmt.Print("- Confirm password: ")
	pass, _ = gopass.GetPasswdMasked()
	repwd := string(pass)
	if pwd != repwd {
		color.Yellow("Password and confirmation are different, please try again")
		fmt.Println("")
		return getConfirmedPasswords()
	}
	return pwd
}

func adminLogin(m *manager.Manager, adminUsr, adminPass string) (token string) {
	var err error

	if adminUsr == "" || adminPass == "" {
		fmt.Println("")
		fmt.Println("Please login with admin credentials:")
	}

	if adminUsr == "" {
		adminUsr = reAsk("- Admin user:")
	}

	if adminPass == "" {
		fmt.Printf("- Admin password: ")
		pass, _ := gopass.GetPasswdMasked()
		adminPass = string(pass)
	}

	if token, err = m.Login(adminUsr, adminPass, ""); err != nil {
		fmt.Println("Invalid credentials, please try again")
		return adminLogin(m, "", "")
	}
	color.Green("Welcome " + adminUsr + "!")

	return token

}

func reAsk(req string) (res string) {
	fmt.Print(req + " ")
	if _, err := fmt.Scanf("%s", &res); err != nil {
		return reAsk(req)
	}
	return res
}

func reAskDef(req, def string) (res string) {
	fmt.Print(req + " ")
	if _, err := fmt.Scanf("%s", &res); err != nil {
		res = def
	}

	return res
}
