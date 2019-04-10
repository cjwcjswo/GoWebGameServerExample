package handler

import (
	"GoWebGameServerExample/pkg/log"
	"GoWebGameServerExample/pkg/protocol"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type GatewayRequest struct {
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

func (request *GatewayRequest) Decoding(data []byte) (protocol.ErrorCode, string) {
	err := json.Unmarshal(data, &request)
	if err != nil {
		log.LocalLogger.Error("GatewayRequest Decoding Fail", zap.String("Error", err.Error()))
		return protocol.ERROR_GATEWAY_DECODING_FAIL, err.Error()
	}
	if request.Seq == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Sequence is nil"
	}
	if request.ApiKey == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "ApiKey is nil"
	}
	if request.ApiVer == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "ApiVer is nil"
	}
	if request.Ts == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Ts is nil"
	}
	if request.Uid == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Uid is nil"
	}
	if request.Sid == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Sid is nil"
	}
	if request.Version == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Version is nil"
	}
	if request.Language == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Language is nil"
	}
	if request.Resend == nil {
		return protocol.ERROR_INVALID_REQUEST_VALUE, "Resend is nil"
	}

	return protocol.ERROR_SUCCESS, ""
}

type GatewayResponse struct {
	ApiResult    GatewayApiResult `json:"api_result"`
	ResponseData interface{}      `json:"response_data"`
}

type GatewayApiResult struct {
	Seq       int64              `json:"seq"`
	ApiKey    string             `json:"api_key"`
	ApiVer    int16              `json:"api_ver"`
	Ts        int64              `json:"ts"`
	Uid       int64              `json:"uid"`
	ErrorCode protocol.ErrorCode `json:"code"`
	Message   string             `json:"message"`
}

type GatewayHandler struct{}

func (handler *GatewayHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	var response GatewayResponse
	defer func() {
		responseBytes, err := json.Marshal(response)
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

	var gatewayRequest GatewayRequest
	requestBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.LocalLogger.Error("Gateway Request Message Read All Fail", zap.String("Error", err.Error()))
		response.ApiResult.ErrorCode = protocol.ERROR_GATEWAY_REQUEST_READ_FAIL
		response.ApiResult.Message = protocol.ErrorMessageMap[protocol.ERROR_GATEWAY_REQUEST_READ_FAIL]
		response.ResponseData = err.Error()
		return
	}

	if errorCode, errorMessage := gatewayRequest.Decoding(requestBytes); errorCode != protocol.ERROR_SUCCESS {
		response.ApiResult.ErrorCode = errorCode
		response.ApiResult.Message = protocol.ErrorMessageMap[errorCode]
		response.ResponseData = errorMessage
		return
	}

	response.ApiResult = GatewayApiResult{
		Seq:    *gatewayRequest.Seq,
		ApiKey: *gatewayRequest.ApiKey,
		ApiVer: *gatewayRequest.ApiVer,
		Ts:     *gatewayRequest.Ts,
		Uid:    *gatewayRequest.Uid,
	}

	logicHandler, isExist := gameHandlerMap[*gatewayRequest.ApiKey]
	if isExist == false {
		response.ApiResult.ErrorCode = protocol.ERROR_INVALID_API_KEY
		response.ApiResult.Message = protocol.ErrorMessageMap[protocol.ERROR_INVALID_API_KEY]
		return
	}

	errorCode, result := logicHandler.Handle(writer, request, gatewayRequest.Params)
	response.ApiResult.ErrorCode = errorCode
	response.ApiResult.Message = protocol.ErrorMessageMap[errorCode]
	response.ResponseData = result
}
