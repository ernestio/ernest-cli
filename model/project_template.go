/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// ProjectTemplate ...
type ProjectTemplate struct {
	URL      string `yaml:"vcloud-url"`
	Network  string `yaml:"public-network"`
	Org      string `yaml:"org"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Token    string `yaml:"secret_access_key"`
	Secret   string `yaml:"access_key_id "`
	Region   string `yaml:"region"`
	Fake     bool   `yaml:"fake"`
}
