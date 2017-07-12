package helper

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ernestio/ernest-cli/model"
	"github.com/r3labs/sse"
)

type rawhandler struct {
	stream chan *sse.Event
}

func (h *rawhandler) subscribe() error {
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

			m := model.Message{}

			if err := json.Unmarshal(cleanedInput, &m); err != nil {
				return err
			}

			fmt.Println("[" + m.Subject + "] : " + m.Body)
		}
	}
}
