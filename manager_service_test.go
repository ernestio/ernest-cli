/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestListServices(t *testing.T) {
	convey.Convey("Given I get all services", t, func() {
		server := mockRequest("/api/services/", "GET", 200, `[]`)
		m := Manager{URL: server.URL}
		_, err := m.ListServices("token")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestListBuilds(t *testing.T) {
	convey.Convey("Given a service", t, func() {
		convey.Convey("Given I get all builds", func() {
			server := mockRequest("/api/services/foo/builds/", "GET", 200, `[]`)
			m := Manager{URL: server.URL}
			_, err := m.ListBuilds("foo", "token")
			convey.Convey("Then It does not fail", func() {
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceStatus(t *testing.T) {
	convey.Convey("Given a service", t, func() {
		server := mockRequest("/api/services/foo", "GET", 200, `{}`)
		m := Manager{URL: server.URL}
		_, err := m.ServiceStatus("token", "foo")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceBuildStatus(t *testing.T) {
	convey.Convey("Given a service and a build id", t, func() {
		server := mockRequest("/api/services/foo/builds/1234567890", "GET", 200, `{}`)
		m := Manager{URL: server.URL}
		_, err := m.ServiceBuildStatus("token", "foo", "1234567890")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestResetService(t *testing.T) {
	convey.Convey("Given a service", t, func() {
		server := mockRequest("/api/services/foo/reset/", "POST", 200, ``)
		m := Manager{URL: server.URL}
		err := m.ResetService("foo", "token")
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDestroy(t *testing.T) {
	convey.Convey("Given a service", t, func() {
		server := mockRequest("/api/services/foo", "DELETE", 200, `{}`)
		m := Manager{URL: server.URL}
		err := m.Destroy("token", "foo", false)
		convey.Convey("Then It does not fail", func() {
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
