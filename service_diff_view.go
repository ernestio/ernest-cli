/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"

	"github.com/pmezard/go-difflib/difflib"
)

func printServiceDiff(build1 Service, build2 Service) {

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(build1.Definition),
		B:        difflib.SplitLines(build2.Definition),
		FromFile: "Original : " + build1.ID,
		ToFile:   "Current : " + build2.ID,
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Printf(text)

}
