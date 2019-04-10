package handler

import (
	"GoWebGameServerExample/pkg/protocol"
	"encoding/json"
	"net/http"
)

var gameHandlerMap map[string]GameHandlerInterface

type GameHandlerInterface interface {
	GetApiName() string
	Handle(writer http.ResponseWriter, request *http.Request, params *json.RawMessage) (protocol.ErrorCode, interface{})
}

func InitGameHandler() {
	gameHandlerMap = make(map[string]GameHandlerInterface)
	gameHandlerMap[API_NAME_CHECK_UPDATE] = &CheckUpdateHandler{}
}

func GetGateWayHandler() *GatewayHandler {
	return &GatewayHandler{}
}
