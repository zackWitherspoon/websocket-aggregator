package web_socket

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type WebSocketResponse []struct {
	Event   string `json:"ev"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (resp *WebSocketResponse) DebugResponse() {
	marshalledResponse, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debug(fmt.Sprintf("Respose: %s", marshalledResponse))
}
