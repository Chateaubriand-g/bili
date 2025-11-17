package model

type UserDTO struct {
	ID          uint   `json:"id"`
	UserName    string `json:"username"`
	NickName    string `json:"nickname"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type STSResponse struct {
	EndPoint        string
	BucketName      string
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
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
