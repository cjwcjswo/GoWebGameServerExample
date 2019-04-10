package handler

import (
	"GoWebGameServerExample/pkg/protocol"
	"encoding/json"
	"net/http"
)

const API_NAME_CHECK_UPDATE = "CheckUpdate"

type CheckUpdateHandler struct{}

func (handler *CheckUpdateHandler) Handle(writer http.ResponseWriter, request *http.Request, params *json.RawMessage) (protocol.ErrorCode, interface{}) {

	return protocol.ERROR_SUCCESS, nil
}

func (handler *CheckUpdateHandler) GetApiName() string {
	return API_NAME_CHECK_UPDATE
}
