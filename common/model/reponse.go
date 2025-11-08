package model

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func SuccessResponseWithMsg(msg string, data interface{}) APIResponse {
	return APIResponse{
		Code:    200,
		Message: msg,
		Data:    data,
	}
}

func BadResponse(code int, msg string, data interface{}) APIResponse {
	return APIResponse{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
