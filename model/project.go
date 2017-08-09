/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Project : Representation of a project
type Project struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	GroupID         int      `json:"group_id"`
	GroupName       string   `json:"group_name"`
	Region          string   `json:"region"`
	VCloudURL       string   `json:"vcloud_url"`
	VseURL          string   `json:"vse_url"`
	ExternalNetwork string   `json:"external_network"`
	Username        string   `json:"username"`
	Environments    []string `json:"environments"`
	Roles           []string `json:"roles"`
}
