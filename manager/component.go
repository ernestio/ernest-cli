/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	"encoding/json"
)

// FindComponents ...
func (m *Manager) FindComponents(token, datacenter, component, service string) (components []interface{}, err error) {
	body, _, err := m.doRequest("/api/components/"+component+"/?datacenter="+datacenter+"&service="+service, "GET", []byte(""), token, "")
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal([]byte(body), &components)
	if err != nil {
		return nil, err
	}
	return components, err
}
