/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Action : action
type Action struct {
	ID           int    `json:"id"`
	Type         string `json:"type,omitempty"`
	Status       string `json:"status,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	Options      struct {
		Filters     []string `json:"filters,omitempty"`
		BuildID     string   `json:"build_id,omitempty"`
		Environment string   `json:"environment,omitempty"`
		Resolution  string   `json:"resolution,omitempty"`
	} `json:"options,omitempty"`
}
