/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"

	"github.com/ernestio/ernest-cli/manager"
	"github.com/ernestio/ernest-cli/model"
)

// CmdSetup : Setup an ernest instance
var CmdSetup = cli.Command{
	Name:  "setup",
	Usage: "Use it to setup your ernest instance",
	Description: `This command will help you to setup your ernest instance by:
- [ ] configure ernest-cli target
- [ ] create a plain user
- [ ] create a group
- [ ] link the user to the group
- [ ] login as the newly created user.
- [ ] create a new datacenter (optional)
	`,
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
		group := createGroup(adminToken, m)
		if err := m.GroupAddUser(adminToken, usr, group); err != nil {
			color.Red(err.Error())
			return nil
		}
		// Login as plain user
		if token, err = m.Login(usr, pwd); err != nil {
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
		fmt.Print("Would you like to create a datacenter now? (Y/n): ")

		datacenterCreation := askForConfirmation()
		if datacenterCreation {
			createDatacenter(token, m)
		}
		color.Green("Congratulations, you've successfully set up your ernest istance.")
		fmt.Println("")
		color.Green("What's next?")

		if !datacenterCreation {
			color.Green("- Create a datacenter with command 'ernest-cli datacenter create --help'")
		}
		color.Green("- Apply a service definition with 'ernest-cli service apply your-template.yml'")
		color.Green("- Need some extra documentation? Run 'ernest-cli docs'")

		return nil
	},
}

func createDatacenter(token string, m *manager.Manager) {
	dt := reAsk("- Please specify the provider you want to work with? (vcloud / aws)")
	if dt == "aws" {
		createAWSDatacenter(token, m)
	} else if dt == "vcloud" {
		createVCloudDatacenter(token, m)
	} else {
		fmt.Println("Invalid provider. Valid providers are (vcloud and aws)")
		createDatacenter(token, m)
	}
}

func createAWSDatacenter(token string, m *manager.Manager) {
	fmt.Println("In order to create an AWS datacenter we will need some info: ")
	dt := "aws"
	dname := reAsk("- Name:")
	dregion := reAsk("- Region:")
	dkey := reAsk("- Access key id:")
	dsecret := reAsk("- Secret access key:")
	if body, err := m.CreateAWSDatacenter(token, dname, dt, dregion, dkey, dsecret); err != nil {
		color.Red("ERROR: " + body)
		createAWSDatacenter(token, m)
	}
}

func createVCloudDatacenter(token string, m *manager.Manager) {
	fmt.Println("In order to create a VCloud datacenter we will need some info: ")
	dt := "vcloud"
	dname := reAsk("- Name:")
	dusername := reAsk("- Username:")
	fmt.Print("- Password: ")
	pass, _ := gopass.GetPasswdMasked()
	dpassword := string(pass)
	durl := reAsk("- URL:")
	dnetwork := reAsk("- Network:")
	dvse := reAsk("- VSE url:")

	if body, err := m.CreateVcloudDatacenter(token, dname, dt, dusername, dpassword, durl, dnetwork, dvse); err != nil {
		color.Red("ERROR: " + body)
		createVCloudDatacenter(token, m)
	}
}

func createGroup(adminToken string, m *manager.Manager) string {
	fmt.Println("")
	fmt.Println("To collaborate with your teammates ernest will define a group name")
	group := reAskDef("Group name (default):", "default")
	if err := m.CreateGroup(adminToken, group); err != nil {
		color.Red(err.Error() + ". Please try again")
		return createGroup(adminToken, m)
	}
	color.Green("Group '" + group + "' successfully created")

	return group
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

	if token, err = m.Login(adminUsr, adminPass); err != nil {
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
