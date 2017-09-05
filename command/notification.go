/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdNotification subcommand
import (
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListNotifications ...
var ListNotifications = cli.Command{
	Name:        "list",
	Usage:       h.T("notification.list.usage"),
	ArgsUsage:   h.T("notification.list.args"),
	Description: h.T("notification.list.description"),
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
	Name:        "delete",
	Usage:       h.T("notification.delete.usage"),
	ArgsUsage:   h.T("notification.delete.args"),
	Description: h.T("notification.delete.description"),
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
	Name:        "update",
	Usage:       h.T("notification.update.usage"),
	ArgsUsage:   h.T("notification.update.args"),
	Description: h.T("notification.update.description"),
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
	Name:        "add",
	Usage:       h.T("notification.service.add.usage"),
	ArgsUsage:   h.T("notification.service.add.args"),
	Description: h.T("notification.service.add.description"),
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
	Name:        "remove",
	Usage:       h.T("notification.service.rm.usage"),
	ArgsUsage:   h.T("notification.service.rm.args"),
	Description: h.T("notification.service.rm.description"),
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
	Name:        "create",
	Usage:       h.T("notification.create.usage"),
	ArgsUsage:   h.T("notification.create.args"),
	Description: h.T("notification.create.description"),
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
