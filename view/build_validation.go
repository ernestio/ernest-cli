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

	fmt.Printf("\nTest Summary: %s, %s, %d total\n", fmtpassed(passed), fmtfailed(failed), total)
}

func fmtpassed(i int) string {
	if i < 1 {
		return fmt.Sprintf("%d passed", i)
	} else {
		return color.GreenString("%d passed", i)
	}
}

func fmtfailed(i int) string {
	if i > 0 {
		return color.RedString("%d failed", i)
	} else {
		return fmt.Sprintf("%d failed", i)
	}
}
