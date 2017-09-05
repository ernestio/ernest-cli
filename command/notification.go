/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdNotification subcommand
import (
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListNotifications ...
var ListNotifications = cli.Command{
	Name:      "list",
	Usage:     "List available notifications.",
	ArgsUsage: " ",
	Description: `List available notifications.

   Example:
    $ ernest notification list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		notifications, err := m.ListNotifications(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

		view.PrintNotificationList(notifications)

		return nil
	},
}

// DeleteNotification : Will delete the specified notification
var DeleteNotification = cli.Command{
	Name:  "delete",
	Usage: "Deletes an existing notify.",
	Description: `Deletes an existing notify on the targeted instance of Ernest.

   Example:
    $ ernest notify delete <notify_name>


   Example:
	 $ ernest notify delete my_notify
	`,
	ArgsUsage: "<notify_name>",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify a valid name")
			return nil
		}

		name := c.Args()[0]
		m, cfg := setup(c)
		err := m.DeleteNotification(cfg.Token, name)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Notify " + name + " successfully delete")
		return nil
	},
}

// UpdateNotification : Will update the notification specific fields
var UpdateNotification = cli.Command{
	Name:  "update",
	Usage: "Update a new notify.",
	Description: `Update an existing notify on the targeted instance of Ernest.

   Example:
    $ ernest notify update <notify_name> <provider-details>


   Example:
	 $ ernest notify update my_notify '{"url":"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"}'
	`,
	ArgsUsage: "<notify_name> <notify_config>",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify a valid name")
			return nil
		}

		if len(c.Args()) < 2 {
			color.Red("You should specify a notify config options")
			return nil
		}

		name := c.Args()[0]
		notifyConfig := c.Args()[1]
		m, cfg := setup(c)
		err := m.UpdateNotification(cfg.Token, name, notifyConfig)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Notify " + name + " successfully updated")
		return nil
	},
}

// AddServiceToNotification : Creates a new user
var AddServiceToNotification = cli.Command{
	Name:  "add",
	Usage: "Add environment to an existing notify.",
	Description: `Adds a environment to an existing notify.

   Example:
    $ ernest notify add <project_name> <environment_name> <notify_name>


   Example:
	 $ ernest notify add my_project my_env my_notify 
	`,
	ArgsUsage: "<env_name> <notify_name>",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify a valid project name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify a valid environment name")
			return nil
		}
		if len(c.Args()) < 3 {
			color.Red("You should specify a valid notify name")
			return nil
		}

		service := c.Args()[0] + "/" + c.Args()[1]
		notify := c.Args()[2]
		m, cfg := setup(c)
		err := m.AddServiceToNotification(cfg.Token, service, notify, false)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Environment " + service + " successfully attached to " + notify + " notify")
		return nil
	},
}

// RmServiceToNotification : Creates a new user
var RmServiceToNotification = cli.Command{
	Name:  "remove",
	Usage: "Removes an environment to an existing notify.",
	Description: `Removes an environment to an existing notify.

   Example:
    $ ernest notify remove <env_name> <notify_name>


   Example:
	 $ ernest notify remove my_env my_notify 
	`,
	ArgsUsage: "<env_name> <notify_name>",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify a valid environment name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify a valid notify name")
			return nil
		}

		service := c.Args()[0]
		notify := c.Args()[1]
		m, cfg := setup(c)
		err := m.AddServiceToNotification(cfg.Token, service, notify, true)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Environment " + service + " successfully removed from " + notify + " notify")
		return nil
	},
}

// CreateNotification : Creates a new user
var CreateNotification = cli.Command{
	Name:  "create",
	Usage: "Create a new notify.",
	Description: `Create a new notify on the targeted instance of Ernest.

   Example:
    $ ernest notify create <notify_name> <provider_type> <provider-details>


   Example:
	 $ ernest notify create my_notify slack '{"url":"https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"}'
	`,
	ArgsUsage: "<notify_name> <notify_type> <notify_config>",
	Action: func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			color.Red("You should specify a valid name")
			return nil
		}
		if len(c.Args()) < 2 {
			color.Red("You should specify a notify type")
			return nil
		}

		if len(c.Args()) < 3 {
			color.Red("You should specify a notify config options")
			return nil
		}

		name := c.Args()[0]
		notifyType := c.Args()[1]
		notifyConfig := c.Args()[2]
		m, cfg := setup(c)
		_, err := m.CreateNotification(cfg.Token, name, notifyType, notifyConfig)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		color.Green("Notify " + name + " successfully created")
		return nil
	},
}

// CmdNotification ...
var CmdNotification = cli.Command{
	Name:  "notify",
	Usage: "Notification related subcommands",
	Subcommands: []cli.Command{
		ListNotifications,
		CreateNotification,
		UpdateNotification,
		DeleteNotification,
		AddServiceToNotification,
		RmServiceToNotification,
	},
}
