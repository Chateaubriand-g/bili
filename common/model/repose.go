package model

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code: 200,
		Data: data,
	}
}
