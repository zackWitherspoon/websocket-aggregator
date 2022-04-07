package web_socket

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type WebSocketResponse []struct {
	Event   string `json:"ev"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (wsResponse *WebSocketResponse) DebugResponse() {
	marshalledResponse, err := json.MarshalIndent(wsResponse, "", "\t")
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debug("Response: %s", marshalledResponse)
}
