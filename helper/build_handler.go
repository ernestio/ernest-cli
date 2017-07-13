/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ernestio/ernest-cli/model"
	"github.com/gosuri/uilive"
	"github.com/r3labs/sse"
)

type buildhandler struct {
	stream    chan *sse.Event
	writer    *uilive.Writer
	format    string
	failures  []error
	args      []interface{}
	service   model.ServiceEvent
	component model.ComponentEvent
}

func (h *buildhandler) subscribe() error {
	for {
		select {
		case msg, ok := <-h.stream:
			if !ok {
				return nil
			}

			if msg.Data == nil {
				continue
			}

			// clean msg body of any null characters
			cleanedInput := bytes.Trim(msg.Data, "\x00")

			m := make(map[string]interface{})

			err := json.Unmarshal(cleanedInput, &m)
			if err != nil {
				return err
			}

			subject := m["_subject"].(string)

			switch subject {
			case SERVICECREATE, SERVICEDELETE, SERVICEIMPORT:
				h.service = processServiceEvent(m)
				h.format, h.args = renderOutput(h.service)
			case SERVICECREATEDONE, SERVICECREATEERROR, SERVICEDELETEDONE, SERVICEDELETEERROR, SERVICEIMPORTDONE, SERVICEIMPORTERROR:
				h.service = processServiceEvent(m)
				err = renderUpdate(h.service, model.ComponentEvent{}, h.args)
				if err != nil {
					return err
				}
			default:
				h.component = processComponentEvent(m)
				cerr := renderUpdate(model.ServiceEvent{}, h.component, h.args)
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
			case SERVICECREATEDONE, SERVICEDELETEDONE, SERVICEIMPORTDONE:
				return nil
			case SERVICECREATEERROR, SERVICEDELETEERROR, SERVICEIMPORTERROR:
				for _, resourceErr := range h.failures {
					fmt.Printf("Message: %s\n\n", red(resourceErr))
				}
				return errors.New("service task failed with errors")
			}

		}
	}
}
