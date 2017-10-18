/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

// SyncChanges : Pretty print for sync changes
func SyncChanges(body string) {
	var lines []string
	if err := json.Unmarshal([]byte(body), &lines); err != nil {
		color.Red("Unexpected response from ernest")
		return
	}

	if len(lines) == 0 {
		fmt.Println("")
		color.Green("There are no changes associated with this build")
		fmt.Println("")
		return
	}

	color.Red("If rejected, ernest will action the following changes:")
	fmt.Println("")
	for i := range lines {
		fmt.Println(" - " + lines[i])
	}
	fmt.Println("")
}
