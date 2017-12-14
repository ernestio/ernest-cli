/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/manager"
	"github.com/ernestio/ernest-cli/model"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

const (
	// NonEmptuTokenVal ...
	NonEmptuTokenVal = "non_empty_token"
	// NonAdminVal ...
	NonAdminVal = "non_admin"
)

// NoValidation ...
var NoValidation = []string{}

// AuthUsersValidation ...
var AuthUsersValidation = []string{NonEmptuTokenVal}

// NonAdminValidation ...
var NonAdminValidation = []string{NonEmptuTokenVal, NonAdminVal}

type validation func(*manager.Client)

var validations = map[string]validation{
	NonEmptuTokenVal: func(client *manager.Client) {
		if client.Config().Token == "" {
			h.PrintError("You're not allowed to perform this action, please log in")
		}
	},
	NonAdminVal: func(client *manager.Client) {
		session := client.Session().Get()
		if !session.IsAdmin() {
			h.PrintError("You don’t have permissions to perform this action")
		}
	},
}

var session *emodels.Session

// esetup ...
func esetup(c *cli.Context, vals []string) *manager.Client {
	session = nil
	config := model.GetConfig()
	if config == nil {
		config = &model.Config{}
		if c.Command.Name != "target" && c.Command.Name != "setup" {
			h.PrintError("Environment not configured, please use target command")
		}
	}

	client := manager.New(config)
	for _, v := range vals {
		if fn, ok := validations[v]; ok {
			fn(client)
		}
	}

	return client

}

// setup ...
// TODO : Deprecate this to use "esetup"
func setup(c *cli.Context) (*manager.Manager, *model.Config) {
	config := model.GetConfig()
	if config == nil {
		config = &model.Config{}
		if c.Command.Name != "target" && c.Command.Name != "setup" {
			h.PrintError("Environment not configured, please use target command")
		}
	}
	m := manager.Manager{URL: config.URL, Version: c.App.Version}
	return &m, config
}

func stringWithDefault(c *cli.Context, key, def string) (val string) {
	if val = c.String(key); val == "" {
		val = def
	}
	return
}

func paramsLenValidation(c *cli.Context, number int, translationKey string) {
	if len(c.Args()) < number {
		h.PrintError("Please provide required parameters:\n" + h.T(translationKey))
	}
}

func requiredFlags(c *cli.Context, flags []string) {
	errs := []string{}
	for _, flag := range flags {
		if c.String(flag) == "" {
			errs = append(errs, "Please provide a "+flag+" with --"+flag+" flag")
		}
	}
	if len(errs) > 0 {
		h.PrintError(strings.Join(errs, "\n"))
	}
}

func stringFlagND(name, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Usage: usage,
	}
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

type flagDef struct {
	typ   string
	def   interface{}
	mapto string
	req   bool
}

func parseTemplateFlags(c *cli.Context, keys map[string]flagDef) map[string]interface{} {
	var err error
	errs := []string{}
	flags := make(map[string]interface{})
	t := c.String("template")
	if t != "" {
		flags, err = getProjectTemplateAsMap(t)
		h.EvaluateError(err)
	}
	for k, t := range keys {
		mapto := k
		if t.mapto != "" {
			mapto = t.mapto
		}
		if t.def != nil {
			flags[mapto] = t.def
		}
		if t.typ == "string" {
			if c.String(k) != "" {
				flags[mapto] = c.String(k)
			}
		} else {
			if c.Bool(k) != false {
				flags[mapto] = c.Bool(k)
			}
		}
		if t.req {
			if _, ok := flags[mapto]; !ok {
				errs = append(errs, "Specify a valid --"+k+" flag")
			}
		}
	}

	if len(errs) > 0 {
		msgs := []string{"Please, fix the error shown below to continue"}
		for _, e := range errs {
			msgs = append(msgs, "  - "+e)
		}
		h.PrintError(strings.Join(msgs, "\n"))
	}

	return flags
}

func getProjectTemplateAsMap(template string) (map[string]interface{}, error) {
	flags := make(map[string]interface{}, 0)
	payload, err := ioutil.ReadFile(template)
	if err != nil {
		return flags, errors.New("Template file '" + template + "' not found")
	}
	if yaml.Unmarshal(payload, &flags) != nil {
		return flags, errors.New("Template file '" + template + "' is not valid yaml file")
	}
	return flags, nil
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}
