/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/r3labs/diff"
)

type component struct {
	created bool
	deleted bool
	changes []diff.Change
}

type componentlist struct {
	keys       sort.StringSlice
	components map[string]*component
}

func (c *component) add(dc diff.Change) {
	(*c).changes = append((*c).changes, dc)
}

func (cl *componentlist) get(c string) *component {
	if (*cl).components[c] == nil {
		(*cl).keys = append((*cl).keys, c)
		(*cl).components[c] = &component{
			changes: make([]diff.Change, 0),
		}
	}

	return (*cl).components[c]
}

func (cl *componentlist) add(c diff.Change) {
	id := strings.Replace(c.Path[0], "::", ".", 1)
	c.Path = c.Path[1:]

	last := c.Path[len(c.Path)-1]

	x := cl.get(id)

	if last == "_component_id" {
		switch c.Type {
		case diff.CREATE:
			x.created = true
		case diff.DELETE:
			x.deleted = true
		}
		return
	}

	x.add(c)
}

// PrintDiff : prints the diff output from two compared builds
func PrintDiff(c *diff.Changelog) {
	list := componentlist{
		keys:       make(sort.StringSlice, 0),
		components: make(map[string]*component),
	}

	for _, change := range *c {
		list.add(change)
	}

	list.keys.Sort()

	for _, component := range list.keys {
		cmp := list.components[component]

		var header string
		if cmp.created {
			header = fmt.Sprintf("\n + %s", color.GreenString(component))
		} else if cmp.deleted {
			header = fmt.Sprintf("\n - %s", color.RedString(component))
		} else {
			header = fmt.Sprintf("\n ~ %s", color.YellowString(component))
		}

		fmt.Println(header)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetColumnSeparator("")
		table.SetRowLine(false)
		table.SetAutoWrapText(false)

		for _, chng := range cmp.changes {
			if chng.From == nil {
				chng.From = ""
				chng.To = color.GreenString(fmt.Sprint(chng.To))
			} else if chng.To == nil {
				chng.To = ""
				chng.From = color.RedString(fmt.Sprint(chng.From))
			} else {
				chng.From = color.RedString(fmt.Sprint(chng.From))
				chng.To = color.GreenString(fmt.Sprint(chng.To))
			}
			if cmp.created {
				table.Append([]string{color.GreenString("   " + strings.Join(chng.Path, ".")), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
			} else if cmp.deleted {
				table.Append([]string{color.RedString("   " + strings.Join(chng.Path, ".")), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
			} else {
				table.Append([]string{color.YellowString("   " + strings.Join(chng.Path, ".")), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
			}
		}

		table.Render()
	}
}
