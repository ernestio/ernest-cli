package helper

import (
	"encoding/json"

	"github.com/fatih/color"
)

// ResponseError contains fields for handling error responses
type ResponseError struct {
	Message string
}

// ResponseMessage is used to generate an error response
func ResponseMessage(body []byte) ResponseError {
	var e ResponseError
	if err := json.Unmarshal(body, &e); err != nil {
		color.Red(err.Error())
	}
	return e
}
