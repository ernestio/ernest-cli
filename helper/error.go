/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"os"

	"github.com/fatih/color"
)

// PrintError : prints an error and returns
func PrintError(msg string) {
	color.Red(msg)
	os.Exit(1)
}

// EvaluateError : Evaluates an error and exits program
func EvaluateError(err error) {
	if err != nil {
		PrintError(err.Error())
	}
}
