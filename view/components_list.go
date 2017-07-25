/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"encoding/json"
	"fmt"
)

// PrintComponentsList : Pretty print for a components list
func PrintComponentsList(components []interface{}) {
	if len(components) == 0 {
		fmt.Println("There are no components meeting this criteria.")
		return
	}

	for _, component := range components {
		c := component.(map[string]interface{})
		sw := false
		for k, v := range c {
			val, _ := json.Marshal(v)
			if !sw {
				fmt.Println("- " + k + " : " + string(val))
				sw = true
			} else {
				fmt.Println("  " + k + " : " + string(val))
			}
		}
	}
}
