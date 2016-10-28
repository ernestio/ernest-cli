/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
a* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func printLoggerList(loggers []Logger) {
	if len(loggers) == 0 {
		fmt.Println("There are no loggers created yet.")
		return
	}

	for _, l := range loggers {
		if l.Type == "logger" {
			fmt.Println("")
			fmt.Println("Logstash based loggers")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Hostname", "Port", "Timeout", "Region", "Url"})

			port := strconv.Itoa(l.Port)
			timeout := strconv.Itoa(l.Timeout)
			table.Append([]string{l.Hostname, port, timeout})
			table.Render()
		}

		if l.Type == "basic" {
			fmt.Println("")
			fmt.Println("Basic logging file configured on : " + l.Logfile)
		}
	}

}
