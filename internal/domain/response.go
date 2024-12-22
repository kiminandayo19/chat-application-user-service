package domain

type APIBaseResponse struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, status bool, message string, data interface{}) APIBaseResponse {
	return APIBaseResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}
