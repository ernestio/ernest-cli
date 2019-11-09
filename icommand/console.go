/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package icommand

// CmdDocs subcommand
import (
	"flag"
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/ernestio/ernest-cli/command"
	h "github.com/ernestio/ernest-cli/helper"
	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var defaults map[string]string
var userChain []*model.Config
var shell *ishell.Shell

// CmdConsole : Open docs in the default browser
var CmdConsole = cli.Command{
	Name:        "console",
	Usage:       "Interactive ernest shell",
	ArgsUsage:   "console",
	Description: "Interactive ernest shell",
	Action: func(ctx *cli.Context) error {
		client := command.Esetup(ctx, command.AuthUsersValidation)

		// TODO force login if is not logged in
		userChain = append(userChain, client.Config())
		h.Console = true
		shell := ishell.New()
		updatePrompt(shell)
		defaults = map[string]string{}

		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			if recover() != nil {
				shell.Run()
			}
		}()
		// display welcome info.

		shell.Println(" _____ ____  _      _____ ____  _____")
		shell.Println("/  __//  __\\/ \\  /|/  __// ___\\/__ __\\")
		shell.Println("|  \\  |  \\/|| |\\ |||  \\  |    \\  / \\")
		shell.Println("|  /_ |    /| | \\|||  /_ \\___ |  | |")
		shell.Println("\\____\\\\_/\\_\\\\_/  \\|\\____\\\\____/  \\_/")
		shell.Println("---------------------------------------------------")
		shell.Println("")
		shell.Println("Start by typing help")

		// register a function for "greet" command.
		shell.AddCmd(&ishell.Cmd{
			Name: "login",
			Help: "login",
			Func: loginICmd(shell),
		})

		shell.AddCmd(projectICmd(shell, ctx))
		shell.AddCmd(sudoICmd(shell, ctx))
		shell.AddCmd(whoamiICmd(shell, ctx))
		shell.AddCmd(exitICmd(shell, ctx))
		shell.AddCmd(defaultsICmd(shell, ctx))
		shell.AddCmd(envICmd(shell, ctx))
		shell.AddCmd(info(shell, ctx))
		shell.AddCmd(notify(shell, ctx))
		shell.AddCmd(policy(shell, ctx))
		shell.AddCmd(role(shell, ctx))
		shell.AddCmd(schedule(shell, ctx))
		shell.AddCmd(users(shell, ctx))
		shell.AddCmd(logger(shell, ctx))
		shell.AddCmd(logout(shell, ctx))

		// run shell
		shell.Run()

		return nil
	},
}

func execute(cmd cli.Command, ctx *cli.Context) {
	fmt.Println("")
	cmd.Run(ctx)
	fmt.Println("")
}

func mapArgs(c *ishell.Context, inputArgs map[string]input) []string {
	args := []string{""}
	addArgs(c, inputArgs, &args)
	return args
}

func addArgs(c *ishell.Context, inputArgs map[string]input, args *[]string) {
	if len(inputArgs) > 0 {
		c.ShowPrompt(false)
		defer c.ShowPrompt(true) // yes, revert after login.
		for name, v := range inputArgs {
			val := readArg(c, name, v.out)
			*args = append(*args, val)
		}
	}
}

func readArg(c *ishell.Context, name, input string) (val string) {
	if def, ok := defaults[name]; ok && def != "" {
		fmt.Printf("%v %v %v %v\n", color.BlueString("Info:"), " Using default", name, color.GreenString(def))
		val = def
	} else {
		val = readLine(c, input)
	}
	return
}

func mapFlags(c *ishell.Context, inputArgs map[string]input) map[string]string {
	flags := make(map[string]string, 0)
	for k, v := range inputArgs {
		flags[k] = readArg(c, k, v.out)
	}
	return flags
}

func readLine(c *ishell.Context, question string) string {
	c.Print(question)
	val := c.ReadLine()
	return val
}

func setDefaults(c *ishell.Context, input map[string]string) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	for k, v := range input {
		defaults[k] = readLine(c, v)
	}

	c.Println("")
	color.Green("Defaults successfully stored")
	printDefaults(c)
	c.Println("")
}

func printDefaults(c *ishell.Context) {
	for k, v := range defaults {
		c.Println(" -", k, ":", v)
	}
}

type input struct {
	out string
}

type value struct {
	val string
}

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }
func (s *stringValue) String() string   { return string(*s) }
func (s *stringValue) Value() string    { return string(*s) }

func getContext(ctx *cli.Context, args []string, flags map[string]string) *cli.Context {
	set := flag.FlagSet{}
	for k, v := range flags {
		args = append(args, "--"+k+"="+v)
	}

	if err := set.Parse(args); err != nil {
		println(err.Error())
	}

	set.Parsed()
	result := cli.NewContext(ctx.App, &set, nil)
	return result
}

func login(c *ishell.Context) (*model.Config, error) {
	var username string
	if len(c.Args) > 0 {
		username = c.Args[0]
	} else {
		c.Print("Username: ")
		username = c.ReadLine()
	}

	// get password.
	c.Print("Password for " + username + ": ")
	password := c.ReadPassword()

	client := command.Elogin(username, password, "")

	token, err := client.Cli().Authenticate()
	if err != nil {
		return nil, err
	}

	cfg := client.Config()
	cfg.Token = token
	cfg.User = username

	if err := model.SaveConfig(client.Config()); err != nil {
		h.PrintError("Can't write config file")
	}

	return cfg, nil
}

func updatePrompt(shell *ishell.Shell) {
	last := userChain[len(userChain)-1]
	shell.SetPrompt(color.BlueString("$ " + last.User + "> "))
}
