package helper

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
)

func renderUpdate(s model.BuildEvent, c model.ComponentEvent, a []interface{}) error {
	// component status
	for i, v := range a {
		t := formatType(c.Type)
		if v == t || v == t+"s" {
			if c.Action == "create" && c.State == "running" {
				a[i+3] = yellow("Creating")
			} else if c.Action == "create" && c.State == "completed" {
				a[i+1] = a[i+1].(int) + 1
				if a[i+1] == a[i+2] {
					a[i+3] = green("Created")
				}
			}

			if c.Action == "update" && c.State == "running" {
				a[i+3] = yellow("Updating")
			} else if c.Action == "update" && c.State == "completed" {
				a[i+1] = a[i+1].(int) + 1
				if a[i+1] == a[i+2] {
					a[i+3] = green("Updated")
				}
			}

			if c.Action == "find" && c.State == "running" {
				a[i+3] = yellow("Searching")
			} else if c.Action == "find" && c.State == "completed" {
				if len(c.Components) > 0 {
					a[i+1] = len(c.Components)
					a[i+2] = len(c.Components)
					a[i+3] = green("Found")
				} else {
					a[i+3] = yellow("None")
				}
			}

			if c.Action == "delete" && c.State == "running" {
				a[i+3] = yellow("Deleting")
			} else if c.Action == "delete" && c.State == "completed" {
				a[i+1] = a[i+1].(int) + 1
				if a[i+1] == a[i+2] {
					a[i+3] = green("Deleted")
				}
			}

			if c.State == "errored" {
				a[i+3] = red("Error")
				return errors.New(c.Error)
			}
		}
	}

	// overall status
	switch s.Subject {
	case BUILDDELETEDONE:
		a[len(a)-1] = green("Destroyed")
	case BUILDCREATEDONE:
		a[len(a)-1] = green("Applied")
	case BUILDIMPORTDONE:
		a[len(a)-1] = green("Imported")
	case BUILDCREATEERROR, BUILDDELETEERROR, BUILDIMPORTERROR:
		a[len(a)-1] = red("Error")
	}

	return nil
}

func renderOutput(s model.BuildEvent) (string, []interface{}) {
	var blue = color.New(color.FgBlue).SprintFunc()

	f := "\nEnvironment Name: %s\n"
	a := []interface{}{blue(s.Name)}
	f = f + "Build ID: %s\n\n"
	a = append(a, blue(s.ID))

	if len(s.Changes) == 0 {
		f = f + green("No changes detected\n")
	} else {
		changes := ParseChanges(s.Changes)

		keys := []string{}
		for key := range changes {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		longestKey := 0
		for _, k := range keys {
			if len(k) > longestKey {
				longestKey = len(k)
			}
		}

		for _, k := range keys {
			f = f + "%-" + strconv.Itoa(longestKey+1) + "s %3d/%-3d %s\n"
			t := formatType(k)
			if s.Subject == BUILDIMPORT {
				a = append(a, t, 0, 0, "")
			} else {
				a = append(a, t+"s", 0, changes[k], "")
			}
		}
	}

	f = f + "\nStatus: %s\n\n"
	var status string
	switch s.Subject {
	case BUILDCREATE:
		status = yellow("Applying")
	case BUILDDELETE:
		status = yellow("Destroying")
	case BUILDIMPORT:
		status = yellow("Importing")
	default:
		status = yellow("Unknown")
	}
	a = append(a, status)

	return f, a
}

// ParseChanges takes a list of components, counts duplicates (removing them
// from the list) and returns a map with the type and count.
func ParseChanges(c []model.ComponentEvent) map[string]int {
	seen := map[string]int{}
	for _, v := range c {
		seen[v.Type]++
	}
	return seen
}

func formatType(t string) string {
	t = strings.Replace(t, "_", " ", -1)
	return strings.Title(t)
}

func processBuildEvent(s map[string]interface{}) model.BuildEvent {
	var m model.BuildEvent

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &m,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(s)
	if err != nil {
		panic(err)
	}

	return m
}

func processComponentEvent(c map[string]interface{}) model.ComponentEvent {
	var m model.ComponentEvent

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &m,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(c)
	if err != nil {
		panic(err)
	}

	return m
}
