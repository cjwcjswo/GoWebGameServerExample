package handler

import (
	"GoWebGameServerExample/pkg/protocol"
	"net/http"
)

type InterfaceGameHandler interface {
	GetApiName() string
	Handle(writer http.ResponseWriter, request *http.Request, gatewayRequest *gatewayRequest) (interface{}, protocol.ServerError)
}
