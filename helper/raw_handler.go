package helper

import (
	"encoding/json"
	"fmt"

	"github.com/ernestio/ernest-cli/model"
)

type rawhandler struct {
	stream chan []byte
}

func (h *rawhandler) subscribe() error {
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

			fmt.Println("[" + m.Subject + "] : " + m.Body)
		}
	}
}
