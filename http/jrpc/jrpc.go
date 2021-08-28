package jrpc

const (
	ParseErrorCode          = -32700
	InvalidRequestErrorCode = -32600
	MethodNotFoundErrorCode = -32601
	InvalidParamsErrorCode  = -32602
	InternalErrorCode       = -32603

	Version = "2.0"
)

var ErrorMessages = map[int]string{
	ParseErrorCode:          "Parse error",
	InvalidRequestErrorCode: "Invalid Request",
	MethodNotFoundErrorCode: "Method not found",
	InvalidParamsErrorCode:  "Invalid params",
	InternalErrorCode:       "Internal error",
}
