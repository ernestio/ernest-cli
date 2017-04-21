package helper

import (
	"encoding/json"

	"github.com/fatih/color"
)

type ResponseError struct {
	Message string
}

func ResponseMessage(body []byte) ResponseError {
	var e ResponseError
	if err := json.Unmarshal(body, &e); err != nil {
		color.Red(err.Error())
	}
	return e
}
