/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

// Login command
// Login with your Ernest credentials
var Login = cli.Command{
	Name:      "login",
	Aliases:   []string{"l"},
	ArgsUsage: " ",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "User credentials",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Password credentials",
		},
	},
	Usage: "Login with your Ernest credentials.",
	Description: `Logs an user into Ernest instance.

   Example:
    $ ernest login

   It can also be used without asking the username and password.

   Example:
    $ ernest login --user <user> --password <password>
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if m == nil {
			os.Exit(1)
		}
		var username string
		var password string

		if c.String("user") == "" {
			fmt.Printf("Username: ")
			_, err := fmt.Scanf("%s", &username)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			}
		} else {
			username = c.String("user")
		}

		if c.String("password") == "" {
			fmt.Printf("Password: ")
			pass, _ := gopass.GetPasswdMasked()
			password = string(pass)
		} else {
			password = c.String("password")
		}

		token, username, err := m.Login(username, password)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		cfg.Token = token
		cfg.User = username
		err = saveConfig(cfg)
		if err != nil {
			color.Red("Can't write config file")
			os.Exit(1)
		}
		color.Green("Log in succesful.")
		return nil
	},
}
