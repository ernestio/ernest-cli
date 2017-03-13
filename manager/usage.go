/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"errors"
)

// GetUsageReport : Get the usage report
func (m *Manager) GetUsageReport(token, from, to string) (body string, err error) {
	body, resp, err := m.doRequest("/api/reports/usage/?from="+from+"&to="+to, "GET", []byte(""), token, "")
	if err != nil {
		if resp.StatusCode == 403 {
			return body, errors.New("You're not allowed to perform this action, please log in with an admin account")
		}

		return body, errors.New(string(body))
	}

	return body, nil
}
