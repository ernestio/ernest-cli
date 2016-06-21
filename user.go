/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// CmdUser subcommand
import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

// ListUsers ...
var ListUsers = cli.Command{
	Name:      "list",
	Usage:     "List available users.",
	ArgsUsage: " ",
	Description: `List available users.

   Example:
    $ ernest user list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		users, err := m.ListUsers(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		fmt.Fprintln(w, "NAME\tID\tEMAIL")
		for _, user := range users {
			str := fmt.Sprintf("%s\t%s\t%s", user.Name, user.ID, user.Email)
			fmt.Fprintln(w, str)
		}
		w.Flush()
		return nil
	},
}

// CreateUser ...
var CreateUser = cli.Command{
	Name:  "create",
	Usage: "Create a new user.",
	Description: `Create a new user on the targeted instance of Ernest.

   Example:
    $ ernest user create --user <adminuser> --password <adminpassword> <username> <password>

   You can also add an email to the user with the flag --email

   Example:
    $ ernest user create --user <adminuser> --password <adminpassword> --email username@example.com <username> <password>
	`,
	ArgsUsage: "<username> <password>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Admin user credentials",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Admin password credentials",
		},
		cli.StringFlag{
			Name:  "email",
			Value: "",
			Usage: "Email for the user",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify an user username and a password"
			color.Red(msg)
			return errors.New(msg)
		}
		usr := c.Args()[0]
		email := c.String("email")
		pwd := c.Args()[1]
		user := c.String("user")
		if user == "" {
			msg := "Password not specified"
			color.Red(msg)
			return errors.New("Password not specified")
		}
		password := c.String("password")
		if password == "" {
			msg := "Password not specified"
			color.Red(msg)
			return errors.New("Password not specified")
		}
		m, _ := setup(c)
		err := m.CreateUser(usr, email, usr, pwd, user, password)
		if err != nil {
			color.Red(err.Error())
			return err
		}
		return nil
	},
}

// PasswordUser ...
var PasswordUser = cli.Command{
	Name:      "password",
	Usage:     "Change password of available users.",
	ArgsUsage: "<user-id>",
	Description: `Change password of available users.

   Example:
    $ ernest user password <user-id>

    or changing a password by being admin:

    $ ernest user password --user <adminuser> --password <adminpassword> <user-id>
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Admin user credentials",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Admin password credentials",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify an user ID"
			color.Red(msg)
			return errors.New("You should specify an user ID")
		}

		m, cfg := setup(c)
		usr := c.Args()[0]

		adminuser := c.String("user")
		adminpassword := c.String("password")

		session, err := m.getSession(cfg.Token)
		if err != nil {
			color.Red(err.Error())
		}

		cuser, err := m.getUser(cfg.Token, session.UserID)
		if err != nil {
			color.Red(err.Error())
		}

		if adminuser != "" && adminpassword != "" || cuser.IsAdmin {
			token := ""
			if cuser.IsAdmin {
				token = cfg.Token
			} else {
				token, err = m.Login(adminuser, adminpassword)
				if err != nil {
					color.Red(err.Error())
					return err
				}
			}

			user, err := m.GetUser(token, usr)
			if err != nil {
				color.Red(err.Error())
				return err
			}

			fmt.Printf("New Password: ")
			npass, _ := gopass.GetPasswdMasked()
			newpassword := string(npass)

			fmt.Printf("Repeat new Password: ")
			rnpass, _ := gopass.GetPasswdMasked()
			rnewpassword := string(rnpass)

			if newpassword != rnewpassword {
				msg := "New password doesn't match."
				color.Red(msg)
				return err
			}

			err = m.ChangePasswordByAdmin(token, user.ID, newpassword)
			if err != nil {
				color.Red(err.Error())
				return err
			}
		} else {

			user, err := m.GetUser(cfg.Token, usr)
			if err != nil {
				color.Red(err.Error())
				return err
			}

			fmt.Printf("Old Password: ")
			opass, _ := gopass.GetPasswdMasked()
			oldpassword := string(opass)

			fmt.Printf("New Password: ")
			npass, _ := gopass.GetPasswdMasked()
			newpassword := string(npass)

			fmt.Printf("Repeat new Password: ")
			rnpass, _ := gopass.GetPasswdMasked()
			rnewpassword := string(rnpass)

			if newpassword != rnewpassword {
				msg := "New password doesn't match."
				color.Red(msg)
				return errors.New("New password doesn't match.")
			}

			err = m.ChangePassword(cfg.Token, user.ID, oldpassword, newpassword)
			if err != nil {
				color.Red(err.Error())
				return err
			}
		}
		return nil
	},
}

// DisableUser ...
var DisableUser = cli.Command{
	Name:  "disable",
	Usage: "Disable available users.",
	Description: `Disable available users.

	Example:
	 $ ernest user disable --user <adminuser> --password <adminpassword> <user-id>
 `,
	ArgsUsage: "<username>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Admin user credentials",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Admin password credentials",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			msg := "You should specify an user ID"
			color.Red(msg)
			return errors.New("You should specify an user ID")
		}

		m, _ := setup(c)
		usr := c.Args()[0]

		adminuser := c.String("user")
		if adminuser == "" {
			msg := "Password not specified"
			color.Red(msg)
			return errors.New("Password not specified")
		}
		adminpassword := c.String("password")
		if adminpassword == "" {
			msg := "Password not specified"
			color.Red(msg)
			return errors.New("Password not specified")
		}

		token, err := m.Login(adminuser, adminpassword)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		m.ChangePasswordByAdmin(token, usr, randString(16))

		color.Green("User successfully disabled.")
		return nil
	},
}

// CmdUser ...
var CmdUser = cli.Command{
	Name:  "user",
	Usage: "User related subcommands",
	Subcommands: []cli.Command{
		ListUsers,
		CreateUser,
		PasswordUser,
		DisableUser,
	},
}
