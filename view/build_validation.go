/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"

	"github.com/ernestio/ernest-go-sdk/models"
	"github.com/fatih/color"
)

func PrintValidation(v *models.Validation) {
	if v == nil {
		return
	}

	if v.Passed() {
		color.Green("Validation Passed!")
	} else {
		color.Red("Validation Failed!")
	}

	var current string
	passed, failed, total := v.Stats()

	for _, control := range v.Controls {
		if control.PolicyName() != current {
			fmt.Println(control.PolicyName())
		}

		current = control.PolicyName()

		if control.Status == "passed" {
			color.Green("✔ %s", control.CodeDesc)
		} else {
			color.Red("✘ %s", control.CodeDesc)
			color.Red("✘ %s", control.CodeDesc)
		}
	}

	fmt.Println("Test Summary: %d passed, %d failed, %d total", passed, failed, total)
}
