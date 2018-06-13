package helper

import (
	"encoding/json"
	"fmt"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

type loghandler struct {
	stream chan []byte
}

func (h *loghandler) subscribe() error {
	for {
		select {
		case msg, ok := <-h.stream:
			if !ok {
				return nil
			}
			if msg == nil {
				continue
			}

			m := model.Message{}

			if err := json.Unmarshal(msg, &m); err != nil {
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
