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
		tStringFlag("login.flags.user"),
		tStringFlag("login.flags.password"),
		tStringFlag("login.flags.verification"),
	},
	Action: func(c *cli.Context) error {
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

		verificationCode = c.String("verification-code")
		client := elogin(username, password, verificationCode)
		token, err := client.Cli().Authenticate()

		// MFA check
		if err != nil && err.Error() == "mfa required" {
			fmt.Printf("Verification code: ")
			vc, _ := gopass.GetPasswdMasked()
			verificationCode = string(vc)

			client = elogin(username, password, verificationCode)
			token, err = client.Cli().Authenticate()
		}

		if err != nil {
			h.PrintError(err.Error())
		}

		cfg := client.Config()
		cfg.Token = token

		if err := model.SaveConfig(client.Config()); err != nil {
			h.PrintError("Can't write config file")
		}

		color.Green("Welcome back " + username)

		return nil
	},
}
