/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// User ...
type User struct {
	ID        int      `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	GroupID   int      `json:"group_id"`
	GroupName string   `json:"group_name"`
	Type      string   `json:"type"`
	IsAdmin   bool     `json:"admin"`
	Projects  []string `json:"projects"`
	Envs      []string `json:"envs"`
}
