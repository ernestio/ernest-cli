/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"fmt"
	"os"
	"runtime"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

// Login command
// Login with your Ernest credentials
var Login = cli.Command{
	Name:        "login",
	Aliases:     []string{"l"},
	Usage:       h.T("login.usage"),
	ArgsUsage:   h.T("login.args"),
	Description: h.T("login.description"),
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
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if m == nil {
			os.Exit(1)
		}

		var username string
		var password string
		var verificationCode string

		if c.String("user") == "" {
			fmt.Printf("Username: ")
			_, err := fmt.Scanln(&username)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			}
		} else {
			username = c.String("user")
		}

		if c.String("password") == "" {
			fmt.Printf("Password: ")
			if runtime.GOOS == "windows" {
				_, err := fmt.Scanln(&password)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				}
			} else {
				pass, _ := gopass.GetPasswdMasked()
				password = string(pass)
			}
		} else {
			password = c.String("password")
		}

		token, err := m.Login(username, password, "")

		// MFA check
		if err != nil && err.Error() == "mfa required" {
			fmt.Printf("Verification code: ")
			vc, _ := gopass.GetPasswdMasked()
			verificationCode = string(vc)

			token, err = m.Login(username, password, verificationCode)
			if err != nil {
				h.PrintError("Authentication failed")
			}
		}

		if err != nil {
			h.PrintError(err.Error())
		}

		cfg.Token = token
		cfg.User = username

		err = model.SaveConfig(cfg)
		if err != nil {
			h.PrintError("Can't write config file")
		}

		color.Green("Welcome back " + username)

		return nil
	},
}
