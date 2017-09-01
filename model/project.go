/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Project : Representation of a project
type Project struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Credentials  map[string]interface{} `json:"credentials"`
	Environments []string               `json:"environments"`
	Roles        []string               `json:"roles"`
}
