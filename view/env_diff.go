/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"

	"github.com/pmezard/go-difflib/difflib"
)

// PrintEnvDiff : Pretty print for an envs diff
func PrintEnvDiff(id1, id2 string, build1, build2 []byte) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(build1)),
		B:        difflib.SplitLines(string(build2)),
		FromFile: "Original : " + id1,
		ToFile:   "Current : " + id2,
		Context:  3,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		panic(err)
	}

	fmt.Printf(text)
}
