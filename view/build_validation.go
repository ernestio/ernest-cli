/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"
	"strings"

	"github.com/ernestio/ernest-go-sdk/models"
	"github.com/fatih/color"
)

func PrintValidation(v *models.Validation) {
	if v == nil {
		return
	}

	passed, failed, total := v.Stats()

	for i, profile := range v.Profiles {
		if i > 0 {
			fmt.Printf("\nPolicy: %s\n\n", profile.PolicyName())
		} else {
			fmt.Printf("Policy: %s\n\n", profile.PolicyName())
		}

		for i, control := range profile.Controls {
			var nl string
			if i > 0 {
				nl = "\n"
			}

			if control.Passed() {
				color.Green("%s    ✔ %s", nl, control.ControlTitle())
			} else {
				color.Red("%s    ✘ %s", nl, control.ControlTitle())
			}

			for _, result := range control.Results {
				desc := strings.Split(result.CodeDesc, ":: ")
				if result.Status == "passed" {
					color.Green("      ✔ %s", desc[len(desc)-1])
				} else {
					color.Red("      ✘ %s", desc[len(desc)-1])
					color.Red("        %s", result.Message)
				}
			}
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
