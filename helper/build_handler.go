/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ernestio/ernest-cli/model"
	"github.com/gosuri/uilive"
)

type buildhandler struct {
	stream    chan []byte
	writer    *uilive.Writer
	format    string
	failures  []error
	args      []interface{}
	service   model.BuildEvent
	component model.ComponentEvent
}

func (h *buildhandler) subscribe() error {
	for {
		select {
		case msg, ok := <-h.stream:
			if !ok {
				return nil
			}

			if msg == nil {
				continue
			}

			m := make(map[string]interface{})

			err := json.Unmarshal(msg, &m)
			if err != nil {
				return err
			}

			subject := m["_subject"].(string)

			switch subject {
			case BUILDCREATE, BUILDDELETE, BUILDIMPORT:
				h.service = processBuildEvent(m)
				h.format, h.args = renderOutput(h.service)
			case BUILDCREATEDONE, BUILDCREATEERROR, BUILDDELETEDONE, BUILDDELETEERROR, BUILDIMPORTDONE, BUILDIMPORTERROR:
				h.service = processBuildEvent(m)
				err = renderUpdate(h.service, model.ComponentEvent{}, h.args)
				if err != nil {
					return err
				}
			default:
				h.component = processComponentEvent(m)
				cerr := renderUpdate(model.BuildEvent{}, h.component, h.args)
				if cerr != nil {
					h.failures = append(h.failures, cerr)
				}
			}

			fmt.Fprintf(h.writer, h.format, h.args...)

			err = h.writer.Flush()
			if err != nil {
				return err
			}

			switch subject {
			case BUILDCREATEDONE, BUILDDELETEDONE, BUILDIMPORTDONE:
				return nil
			case BUILDCREATEERROR, BUILDDELETEERROR, BUILDIMPORTERROR:
				for _, resourceErr := range h.failures {
					fmt.Printf("Message: %s\n\n", red(resourceErr))
				}
				return errors.New("service task failed with errors")
			}

		}
	}
}
