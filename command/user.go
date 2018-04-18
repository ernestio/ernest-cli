/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdUser subcommand
import (
	"crypto/rand"
	"fmt"
	"math/big"
	"unicode"

	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"

	h "github.com/ernestio/ernest-cli/helper"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// ListUsers : Gets a list of accessible users
var ListUsers = cli.Command{
	Name:        "list",
	Usage:       h.T("user.list.usage"),
	ArgsUsage:   h.T("user.list.args"),
	Description: h.T("user.list.description"),
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		users := client.User().List()
		view.PrintUserList(users)

		return nil
	},
}

// CreateUser : Creates a new user
var CreateUser = cli.Command{
	Name:        "create",
	Usage:       h.T("user.create.usage"),
	ArgsUsage:   h.T("user.create.args"),
	Description: h.T("user.create.description"),
	Flags: []cli.Flag{
		tStringFlag("user.create.flags.email"),
		tBoolFlag("user.create.flags.mfa"),
		tBoolFlag("user.create.flags.admin"),
	},
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "user.create.args")
		client := esetup(c, AuthUsersValidation)
		usr := c.Args()[0]
		mfa := c.Bool("mfa")

		user := &emodels.User{
			Username: usr,
			Email:    c.String("email"),
			Password: c.Args()[1],
			MFA:      &mfa,
			Disabled: false,
		}
		client.User().Create(user)
		color.Green("User %s successfully created\n\n", usr)

		if mfa {
			color.Green("MFA enabled")
			fmt.Printf("Account name: Ernest (%s)\nKey: %s\n", usr, user.MFASecret)
		}

		if c.Bool("admin") {
			client.User().Promote(user)
		}

		return nil
	},
}

// PasswordUser : Allows users or admins to change its passwords
var PasswordUser = cli.Command{
	Name:        "change-password",
	Usage:       h.T("user.change_password.usage"),
	Description: h.T("user.change_password.description"),
	Flags: []cli.Flag{
		tStringFlag("user.change_password.flags.user"),
		tStringFlag("user.change_password.flags.password"),
		tStringFlag("user.change_password.flags.current-password"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		session := client.Session().Get()

		username := c.String("user")
		password := c.String("password")
		currentPassword := c.String("current-password")

		if !session.IsAdmin() && username != "" {
			h.PrintError("You don’t have permissions to perform this action")
		}

		if session.IsAdmin() && username != "" {

			if password == "" {
				h.PrintError("Please provide a valid password for the user with `--password`")
			}
			user := client.User().Get(username)
			user.Password = password
			user.Disabled = false
			client.User().Update(user)
			color.Green("`" + username + "` password has been changed")

		} else {

			// Ask the user for credentials
			users := client.User().List()
			if len(users) == 0 {
				h.PrintError("You don’t have permissions to perform this action")
			}

			oldpassword := currentPassword
			newpassword := password
			rnewpassword := password

			if oldpassword == "" || newpassword == "" {
				fmt.Printf("You're about to change your password, please respond the questions below: \n")
				fmt.Printf("Current password: ")
				opass, _ := gopass.GetPasswdMasked()
				oldpassword = string(opass)

				fmt.Printf("New password: ")
				npass, _ := gopass.GetPasswdMasked()
				newpassword = string(npass)

				fmt.Printf("Confirm new password: ")
				rnpass, _ := gopass.GetPasswdMasked()
				rnewpassword = string(rnpass)
			}

			if newpassword != rnewpassword {
				h.PrintError("Aborting... New password and confirmation doesn't match.")
			}

			username := client.Config().User
			user := client.User().Get(username)
			user.Password = newpassword
			user.OldPassword = oldpassword
			user.Disabled = false
			client.User().Update(user)

			color.Green("Your password has been changed")
		}

		return nil
	},
}

// DisableUser : Will disable a user (change its password)
var DisableUser = cli.Command{
	Name:        "disable",
	Usage:       h.T("user.disable.usage"),
	ArgsUsage:   h.T("user.disable.args"),
	Description: h.T("user.disable.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "user.disable.args")
		client := esetup(c, NonAdminValidation)
		username := c.Args()[0]

		user := client.User().Get(username)
		user.Password = randString(16)
		user.Disabled = true
		client.User().Update(user)

		color.Green("Account `" + username + "` has been disabled")
		return nil
	},
}

// InfoUser :
var InfoUser = cli.Command{
	Name:        "info",
	Usage:       h.T("user.info.usage"),
	Description: h.T("user.info.description"),
	Flags: []cli.Flag{
		tStringFlag("user.info.flags.user"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, NonAdminValidation)
		username := stringWithDefault(c, "user", client.Config().User)

		user := client.User().Get(username)
		view.PrintUserInfo(user)
		return nil
	},
}

// AddAdminUser :
var AddAdminUser = cli.Command{
	Name:        "add",
	Usage:       h.T("user.admin.add.usage"),
	ArgsUsage:   h.T("user.admin.add.args"),
	Description: h.T("user.admin.add.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "user.admin.add.args")

		client := esetup(c, NonAdminValidation)
		username := c.Args()[0]

		user := client.User().Get(username)
		user.Admin = true
		client.User().Update(user)

		color.Green("Admin privileges assigned to user " + username)
		return nil
	},
}

// RmAdminUser :
var RmAdminUser = cli.Command{
	Name:        "rm",
	Usage:       h.T("user.admin.rm.usage"),
	Description: h.T("user.admin.rm.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 1, "user.admin.rm.args")

		client := esetup(c, NonAdminValidation)
		username := c.Args()[0]

		user := client.User().Get(username)
		user.Admin = false
		client.User().Update(user)

		color.Green("Admin privileges revoked from user " + username)
		return nil
	},
}

// EnableMFA turns on Multi-Factor authentication
var EnableMFA = cli.Command{
	Name:        "enable-mfa",
	Usage:       h.T("user.enable-mfa.usage"),
	Description: h.T("user.enable-mfa.description"),
	Flags: []cli.Flag{
		tStringFlagND("user.enable-mfa.flags.user"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, NonAdminValidation)
		session := client.Session().Get()
		username := stringWithDefault(c, "user-name", session.Username)

		user := client.User().Get(username)
		if user.MFA != nil && *user.MFA {
			fmt.Println("MFA already enabled")
			return nil
		}

		secret := client.User().ToggleMFA(user, true)
		color.Green("MFA enabled")
		fmt.Printf("Account name: Ernest (%s)\nKey: %s\n", user.Username, secret)

		return nil
	},
}

// DisableMFA turns off Multi-Factor authentication
var DisableMFA = cli.Command{
	Name:        "disable-mfa",
	Usage:       h.T("user.disable-mfa.usage"),
	Description: h.T("user.disable-mfa.description"),
	Flags: []cli.Flag{
		tStringFlagND("user.disable-mfa.flags.user"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, NonAdminValidation)
		session := client.Session().Get()
		username := stringWithDefault(c, "user-name", session.Username)

		user := client.User().Get(username)
		if user.MFA == nil || !*user.MFA {
			fmt.Println("MFA already disabled")
			return nil
		}

		_ = client.User().ToggleMFA(user, false)
		color.Red("MFA disabled")

		return nil
	},
}

// ResetMFA generates a new secret for Multi-Factor authentication
var ResetMFA = cli.Command{
	Name:        "reset-mfa",
	Usage:       h.T("user.reset-mfa.usage"),
	Description: h.T("user.reset-mfa.description"),
	Flags: []cli.Flag{
		tStringFlagND("user.reset-mfa.flags.user"),
	},
	Action: func(c *cli.Context) error {
		client := esetup(c, NonAdminValidation)
		session := client.Session().Get()
		username := stringWithDefault(c, "user-name", session.Username)
		user := client.User().Get(username)

		_ = client.User().ToggleMFA(user, false)
		secret := client.User().ToggleMFA(user, true)

		color.Green("MFA reset")
		fmt.Printf("Account name: Ernest (%s)\nKey: %s\n", user.Username, secret)

		return nil
	},
}

// generate random string
func randString(n int) string {
	g := big.NewInt(0)
	max := big.NewInt(130)
	bs := make([]byte, n)

	for i := range bs {
		g, _ = rand.Int(rand.Reader, max)
		r := rune(g.Int64())
		for !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			g, _ = rand.Int(rand.Reader, max)
			r = rune(g.Int64())
		}
		bs[i] = byte(g.Int64())
	}
	return string(bs)
}

// AdminUser ...
var AdminUser = cli.Command{
	Name:  "admin",
	Usage: h.T("user.admin.usage"),
	Subcommands: []cli.Command{
		AddAdminUser,
		RmAdminUser,
	},
}

// CmdUser ...
var CmdUser = cli.Command{
	Name:  "user",
	Usage: h.T("user.usage"),
	Subcommands: []cli.Command{
		CreateUser,
		DisableUser,
		DisableMFA,
		EnableMFA,
		InfoUser,
		ListUsers,
		PasswordUser,
		ResetMFA,
		AdminUser,
	},
}
