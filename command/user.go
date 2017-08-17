/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdUser subcommand
import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/tabwriter"
	"unicode"

	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// ListUsers : Gets a list of accessible users
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
			return nil
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Group", "Admin"})
		for _, u := range users {
			id := strconv.Itoa(u.ID)
			admin := "no"
			if u.IsAdmin {
				admin = "yes"
			}
			table.Append([]string{id, u.Username, u.GroupName, admin})
		}
		table.Render()

		return nil
	},
}

// CreateUser : Creates a new user
var CreateUser = cli.Command{
	Name:  "create",
	Usage: "Create a new user.",
	Description: `Create a new user on the targeted instance of Ernest.

   Example:
    $ ernest user create <username> <password>

   You can also add an email to the user with the flag --email

   Example:
    $ ernest user create --email username@example.com <username> <password>
	`,
	ArgsUsage: "<username> <password>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "email",
			Value: "",
			Usage: "Email for the user",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify an user username and a password")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify the user password")
			return nil
		}

		usr := c.Args()[0]
		email := c.String("email")
		pwd := c.Args()[1]
		m, cfg := setup(c)
		err := m.CreateUser(cfg.Token, usr, email, usr, pwd)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("User " + usr + " successfully created")
		return nil
	},
}

// PasswordUser : Allows users or admins to change its passwords
var PasswordUser = cli.Command{
	Name:  "change-password",
	Usage: "Change password of available users.",
	Description: `Change password of available users.

   Example:
    $ ernest user change-password

    or changing a change-password by being admin:

    $ ernest user change-password --user <username> --current-password <current-password> --password <new-password>
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "The username of the user to change password",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "The new user password",
		},
		cli.StringFlag{
			Name:  "current-password",
			Value: "",
			Usage: "The current user password",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)

		username := c.String("user")
		password := c.String("password")
		currentPassword := c.String("current-password")

		session, err := m.GetSession(cfg.Token)
		if err != nil {
			color.Red("You don’t have permissions to perform this action")
			return nil
		}

		if !session.IsAdmin && username != "" {
			color.Red("You don’t have permissions to perform this action")
			return nil
		}

		if session.IsAdmin && username != "" {
			if password == "" {
				color.Red("Please provide a valid password for the user with `--password`")
				return nil
			}

			// Just change the password with the given values for the given user
			usr, err := m.GetUserByUsername(cfg.Token, username)
			if err = m.ChangePasswordByAdmin(cfg.Token, usr.ID, usr.Username, usr.GroupID, password); err != nil {
				color.Red(err.Error())
				return nil
			}
			color.Green("`" + usr.Username + "` password has been changed")
		} else {
			// Ask the user for credentials
			var users []model.User
			if users, err = m.ListUsers(cfg.Token); err != nil {
				color.Red("You don’t have permissions to perform this action")
				return nil
			}
			if len(users) == 0 {
				color.Red("You don’t have permissions to perform this action")
				return nil
			}

			var user model.User
			for _, u := range users {
				if u.Username == cfg.User {
					user = u
					break
				}
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
				color.Red("Aborting... New password and confirmation doesn't match.")
				return nil
			}

			err = m.ChangePassword(cfg.Token, user.ID, user.Username, user.GroupID, oldpassword, newpassword)
			if err != nil {
				color.Red(err.Error())
				return err
			}
			color.Green("Your password has been changed")
		}

		return nil
	},
}

// DisableUser : Will disable a user (change its password)
var DisableUser = cli.Command{
	Name:  "disable",
	Usage: "Disable available users.",
	Description: `Disable available users.

	Example:
	 $ ernest user disable <user-name>
 `,
	ArgsUsage: "<username>",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify an username")
			return nil
		}

		m, cfg := setup(c)
		username := c.Args()[0]

		session, err := m.GetSession(cfg.Token)
		if err != nil {
			color.Red("You don’t have permissions to perform this action")
			return nil
		}

		if !session.IsAdmin {
			color.Red("You don’t have permissions to perform this action")
			return nil
		}

		user, err := m.GetUserByUsername(cfg.Token, username)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		if err = m.ChangePasswordByAdmin(cfg.Token, user.ID, user.Username, user.GroupID, randString(16)); err != nil {
			color.Red(err.Error())
			return nil
		}

		color.Green("Account `" + username + "` has been disabled")
		return nil
	},
}

// InfoUser :
var InfoUser = cli.Command{
	Name:  "info",
	Usage: "Displays information about the specified user (current user by default).",
	Description: `

	Example:
	 $ ernest user info
	 $ ernest user info --user <user-name>
 `,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Username",
		},
	},
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		session, err := m.GetSession(cfg.Token)
		if err != nil {
			color.Red("You don’t have permissions to perform this action")
			return nil
		}

		username := c.String("user")
		if username != "" && session.IsAdmin == false {
			color.Red("You don’t have permissions to access '" + username + "' information")
			return nil
		}
		if username == "" {
			username = cfg.User
		}

		user, err := m.GetUser(cfg.Token, username)
		if err != nil {
			color.Red(err.Error())
			return err
		}

		view.PrintUserInfo(user)
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

// CmdUser ...
var CmdUser = cli.Command{
	Name:  "user",
	Usage: "User related subcommands",
	Subcommands: []cli.Command{
		ListUsers,
		CreateUser,
		PasswordUser,
		DisableUser,
		InfoUser,
	},
}
