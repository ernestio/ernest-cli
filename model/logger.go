/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// Logger : Representation of a logger
type Logger struct {
	Type        string `json:"type"`
	Logfile     string `json:"logfile"`
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
	Timeout     int    `json:"timeout"`
	Token       string `json:"token"`
	Environment string `json:"environment"`
}
