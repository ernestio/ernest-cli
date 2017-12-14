/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package view

import (
	"fmt"

	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// PrintProjectInfo : Pretty print for a project
func PrintProjectInfo(project *emodels.Project) {
	fmt.Println("Name: ", project.Name)
	fmt.Println("Provider: ")
	fmt.Println("  Type: ", project.Type)
	fmt.Println("Environments: ")
	for _, v := range project.Environments {
		fmt.Println("  ", v)
	}
	fmt.Println("Members: ")
	for _, v := range project.Roles {
		fmt.Println("  ", v)
	}

}
