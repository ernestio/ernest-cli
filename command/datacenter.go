/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

// CmdDatacenter subcommand
import (
	"errors"
	"io/ioutil"

	"github.com/ernestio/ernest-cli/model"
	"github.com/ernestio/ernest-cli/view"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

// ListDatacenters ...
var ListDatacenters = cli.Command{
	Name:      "list",
	Usage:     "List available datacenters.",
	ArgsUsage: " ",
	Description: `List available datacenters.

   Example:
    $ ernest datacenter list
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}
		datacenters, err := m.ListDatacenters(cfg.Token)
		if err != nil {
			color.Red(err.Error())
			return nil
		}

		view.PrintDatacenterList(datacenters)

		return nil
	},
}

// UpdateDatacenters : Will update the datacenter specific fields
var UpdateDatacenters = cli.Command{
	Name:        "update",
	Usage:       "Updates an existing datacenter.",
	Description: "Update an existing datacenter on the targeted instance of Ernest.",
	Subcommands: []cli.Command{
		UpdateVCloudDatacenter,
		UpdateAWSDatacenter,
		UpdateAzureDatacenter,
	},
}

// CreateDatacenters ...
var CreateDatacenters = cli.Command{
	Name:        "create",
	Usage:       "Create a new datacenter.",
	Description: "Create a new datacenter on the targeted instance of Ernest.",
	Subcommands: []cli.Command{
		CreateVcloudDatacenter,
		CreateAWSDatacenter,
		CreateAzureDatacenter,
	},
}

// CmdDatacenter ...
var CmdDatacenter = cli.Command{
	Name:  "datacenter",
	Usage: "Datacenter related subcommands",
	Subcommands: []cli.Command{
		ListDatacenters,
		CreateDatacenters,
		UpdateDatacenters,
		DeleteDatacenter,
	},
}

func getDatacenterTemplate(template string, t *model.DatacenterTemplate) (err error) {
	payload, err := ioutil.ReadFile(template)
	if err != nil {
		return errors.New("Template file '" + template + "' not found")
	}
	if yaml.Unmarshal(payload, &t) != nil {
		return errors.New("Template file '" + template + "' is not valid yaml file")
	}
	return err
}
