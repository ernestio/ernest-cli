package helper

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/r3labs/sse"
)

type loghandler struct {
	stream chan *sse.Event
}

func (h *loghandler) subscribe() error {
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

			color.Yellow(m.Subject)
			if len(m.Body) > 0 {
				message, _ := prettyjson.Format([]byte(m.Body))
				fmt.Println(string(message))
			} else {
				fmt.Println("-- Empty string --")
			}

		}
	}
}
