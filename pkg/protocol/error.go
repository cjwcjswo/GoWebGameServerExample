package protocol

type ErrorCode int16

type ServerError struct {
	ErrorCode ErrorCode
}

func (error ServerError) Error() string {
	if errorMessage, isExist := ErrorMessageMap[error.ErrorCode]; isExist {
		return errorMessage
	}
	return ""
}

const (
	ERROR_SUCCESS = 200

	ERROR_GATEWAY_REQUEST_READ_FAIL      = 1001
	ERROR_GATEWAY_DECODING_FAIL          = 1002
	ERROR_INVALID_REQUEST_VALUE          = 1003
	ERROR_INVALID_API_KEY                = 1004
	ERROR_MYSQL_FAIL                     = 1005
	ERROR_PARAMS_DECODING_FAIL           = 1006
	ERROR_YOU_NEED_TO_UPDATE_APP_VERSION = 1007
	ERROR_WRONG_PLATFORM                 = 1008
	ERROR_APP_VERSION_DATA_NOT_FOUND     = 1009
)

var ErrorMessageMap = map[ErrorCode]string{
	ERROR_SUCCESS: "Success",

	ERROR_GATEWAY_REQUEST_READ_FAIL:      "Request read fail",
	ERROR_GATEWAY_DECODING_FAIL:          "Gateway decoding fail",
	ERROR_INVALID_REQUEST_VALUE:          "Invalid request value",
	ERROR_INVALID_API_KEY:                "Invalid api key",
	ERROR_MYSQL_FAIL:                     "MySQL Func Fail",
	ERROR_PARAMS_DECODING_FAIL:           "Params Decoding Fail",
	ERROR_YOU_NEED_TO_UPDATE_APP_VERSION: "You Need to Update App Version!",
	ERROR_WRONG_PLATFORM:                 "Wrong Platform!",
	ERROR_APP_VERSION_DATA_NOT_FOUND:     "App Version Data Not Found",
}
