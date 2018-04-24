/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdNotification subcommand
import (
	"fmt"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// ListNotifications ...
var ListNotifications = cli.Command{
	Name:        "list",
	Usage:       h.T("notification.list.usage"),
	ArgsUsage:   h.T("notification.list.args"),
	Description: h.T("notification.list.description"),
	Action: func(c *cli.Context) error {
		client := esetup(c, AuthUsersValidation)
		notifications := client.Notification().List()

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
		paramsLenValidation(c, 1, "notification.delete.args")
		name := c.Args()[0]
		client := esetup(c, AuthUsersValidation)

		client.Notification().Delete(name)
		color.Green(fmt.Sprintf(h.T("notification.delete.success"), name))
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
		paramsLenValidation(c, 2, "notification.update.args")
		client := esetup(c, AuthUsersValidation)
		name := c.Args()[0]
		notificationConfig := c.Args()[1]

		n := client.Notification().Get(name)
		n.Config = notificationConfig
		client.Notification().Update(n)
		color.Green(fmt.Sprintf(h.T("notification.update.success"), name))
		return nil
	},
}

// AddEntityToNotification : Creates a new user
var AddEntityToNotification = cli.Command{
	Name:        "add",
	Usage:       h.T("notification.service.add.usage"),
	ArgsUsage:   h.T("notification.service.add.args"),
	Description: h.T("notification.service.add.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "notification.service.add.args")

		notification := c.Args()[0]
		project := c.Args()[1]
		client := esetup(c, AuthUsersValidation)
		entity := project
		if len(c.Args()) > 2 {
			entity = entity + "/" + c.Args()[2]
			client.Notification().AddEnv(notification, project, c.Args()[2])
		} else {
			client.Notification().AddProject(notification, project)
		}

		color.Green(fmt.Sprintf(h.T("notification.service.add.success"), entity, notification))
		return nil
	},
}

// RmEntityToNotification : Creates a new user
var RmEntityToNotification = cli.Command{
	Name:        "remove",
	Usage:       h.T("notification.service.rm.usage"),
	ArgsUsage:   h.T("notification.service.rm.args"),
	Description: h.T("notification.service.rm.description"),
	Action: func(c *cli.Context) error {
		paramsLenValidation(c, 2, "notification.service.rm.args")

		notification := c.Args()[0]
		project := c.Args()[1]
		client := esetup(c, AuthUsersValidation)
		entity := project
		if len(c.Args()) > 2 {
			entity = entity + "/" + c.Args()[2]
			client.Notification().RmEnv(notification, project, c.Args()[2])
		} else {
			client.Notification().RmProject(notification, project)
		}

		color.Green(fmt.Sprintf(h.T("notification.service.rm.success"), entity, notification))
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
		paramsLenValidation(c, 3, "notification.create.args")
		client := esetup(c, AuthUsersValidation)

		name := c.Args()[0]
		notificationType := c.Args()[1]
		notificationConfig := c.Args()[2]
		notification := emodels.Notification{
			Name:   name,
			Type:   notificationType,
			Config: notificationConfig,
		}
		client.Notification().Create(&notification)
		color.Green(fmt.Sprintf(h.T("notification.create.success"), name))
		return nil
	},
}

// CmdNotification ...
var CmdNotification = cli.Command{
	Name:  "notification",
	Usage: "Notification related subcommands",
	Subcommands: []cli.Command{
		ListNotifications,
		CreateNotification,
		UpdateNotification,
		DeleteNotification,
		AddEntityToNotification,
		RmEntityToNotification,
	},
}
