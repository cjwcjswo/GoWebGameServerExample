package handler

import (
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/protocol"
	"GoWebGameServerExample/pkg/util"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type gatewayRequest struct {
	Seq      *int64           `json:"seq"`
	ApiKey   *string          `json:"api_key"`
	ApiVer   *int16           `json:"api_ver"`
	Ts       *int64           `json:"ts"`
	Uid      *int64           `json:"uid"`
	Sid      *string          `json:"sid"`
	Version  *string          `json:"version"`
	Language *string          `json:"language"`
	Resend   *bool            `json:"re_send"`
	Params   *json.RawMessage `json:"params"`
}

type gatewayResponse struct {
	ApiResult    gatewayApiResult `json:"api_result"`
	ResponseData interface{}      `json:"response_data"`
}

type gatewayApiResult struct {
	Seq       int64              `json:"seq"`
	ApiKey    string             `json:"api_key"`
	ApiVer    int16              `json:"api_ver"`
	Ts        int64              `json:"ts"`
	Uid       int64              `json:"uid"`
	ErrorCode protocol.ErrorCode `json:"code"`
	Message   string             `json:"message"`
}

type GatewayHandler struct {
	handlerMap map[string]InterfaceGameHandler
	parser     util.Parser
}

func NewGatewayHandler() *GatewayHandler {
	handler := new(GatewayHandler)
	handler.parser = new(util.JsonParser) // Parser Setting
	handler.handlerMap = make(map[string]InterfaceGameHandler)

	// Init Logic Handler
	handler.handlerMap[API_NAME_CHECK_UPDATE] = NewCheckUpdateHandler()

	return handler
}

func (handler *GatewayHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	var response gatewayResponse
	defer func() {
		responseBytes, err := handler.parser.Encoding(response)
		if err != nil {
			log.LocalLogger.Error("Response Bytes Marshal Fail", zap.String("Error", err.Error()))
			return
		}
		_, err = writer.Write(responseBytes)
		if err != nil {
			log.LocalLogger.Error("Response Write Fail!", zap.String("Response", string(responseBytes)), zap.String("Error", err.Error()))
			return
		}
	}()

	requestBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.LocalLogger.Error("Gateway Request Message Read All Fail", zap.String("Error", err.Error()))
		response.ApiResult.ErrorCode = protocol.ERROR_GATEWAY_REQUEST_READ_FAIL
		response.ApiResult.Message = protocol.ErrorMessageMap[protocol.ERROR_GATEWAY_REQUEST_READ_FAIL]
		response.ResponseData = err.Error()
		return
	}

	gatewayRequest, serverError := handler.decodingRequest(requestBytes)
	if serverError.ErrorCode != protocol.ERROR_SUCCESS {
		response.ApiResult.ErrorCode = serverError.ErrorCode
		response.ApiResult.Message = serverError.Error()
		return
	}

	response.ApiResult = gatewayApiResult{
		Seq:    *gatewayRequest.Seq,
		ApiKey: *gatewayRequest.ApiKey,
		ApiVer: *gatewayRequest.ApiVer,
		Ts:     *gatewayRequest.Ts,
		Uid:    *gatewayRequest.Uid,
	}

	logicHandler, isExist := handler.handlerMap[*gatewayRequest.ApiKey]
	if isExist == false {
		response.ApiResult.ErrorCode = protocol.ERROR_INVALID_API_KEY
		response.ApiResult.Message = protocol.ErrorMessageMap[protocol.ERROR_INVALID_API_KEY]
		return
	}

	result, serverError := logicHandler.Handle(writer, request, &gatewayRequest)
	response.ApiResult.ErrorCode = serverError.ErrorCode
	response.ApiResult.Message = serverError.Error()
	response.ResponseData = result
}

func (handler *GatewayHandler) decodingRequest(data []byte) (gatewayRequest, protocol.ServerError) {
	var request gatewayRequest
	err := handler.parser.Decoding(data, &request)
	if err != nil {
		log.LocalLogger.Error("GatewayRequest Decoding Fail", zap.String("Error", err.Error()))
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_GATEWAY_DECODING_FAIL}
	}
	if request.Seq == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.ApiKey == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.ApiVer == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.Ts == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.Uid == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.Sid == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.Version == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.Language == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}
	if request.Resend == nil {
		return request, protocol.ServerError{ErrorCode: protocol.ERROR_INVALID_REQUEST_VALUE}
	}

	return request, protocol.ServerError{ErrorCode: protocol.ERROR_SUCCESS}
}
