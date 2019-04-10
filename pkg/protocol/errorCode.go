package protocol

type ErrorCode int16

const (
	ERROR_SUCCESS = 200

	ERROR_GATEWAY_REQUEST_READ_FAIL = 1001
	ERROR_GATEWAY_DECODING_FAIL     = 1002
	ERROR_INVALID_REQUEST_VALUE     = 1003
	ERROR_INVALID_API_KEY           = 1004
)

var ErrorMessageMap = map[ErrorCode]string{
	ERROR_SUCCESS: "Success",

	ERROR_GATEWAY_REQUEST_READ_FAIL: "Request read fail",
	ERROR_GATEWAY_DECODING_FAIL:     "Gateway decoding fail",
	ERROR_INVALID_REQUEST_VALUE:     "Invalid request value",
	ERROR_INVALID_API_KEY:           "Invalid api key",
}
