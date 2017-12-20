/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
a* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintLoggerList : pretty print for loggers list
func PrintLoggerList(loggers []*emodels.Logger) {
	if len(loggers) == 0 {
		fmt.Println("There are no loggers created yet.")
		return
	}

	for _, l := range loggers {
		if l.Type == "logstash" {
			fmt.Println("")
			fmt.Println("Logstash based loggers")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Hostname", "Port", "Timeout"})

			port := strconv.Itoa(l.Port)
			timeout := strconv.Itoa(l.Timeout)
			table.Append([]string{l.Hostname, port, timeout})
			table.Render()
		}

		if l.Type == "basic" {
			fmt.Println("")
			fmt.Println("Your basic logfile is configured on : " + l.Logfile)
			fmt.Println("* In case you're running dockerized ernest the file refers to the container path, you'll find your mapped volume on docker-compose.yml file of your ernest usually ernest/logs")
		}

		if l.Type == "rollbar" {
			fmt.Println("")
			fmt.Println("Rollbar logging is active")
		}
	}

}
