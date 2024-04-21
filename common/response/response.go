// Package response ...
package response

import "net/http"

// BaseResponse reusable response
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse API success response
type SuccessResponse struct {
	*BaseResponse
	Data interface{} `json:"data"`
}

var OkMessage = SuccessResponse{
	BaseResponse: &BaseResponse{
		Code:    http.StatusOK,
		Message: "ok",
	},
}

// ErrorResponse API error response
type ErrorResponse struct {
	*BaseResponse
	Detail string `json:"detail"`
}
