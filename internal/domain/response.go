package domain

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"totalData"`
	TotalPage int   `json:"totalPage"`
}

type APIBaseResponse struct {
	Code       int         `json:"code"`
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"` // Changed to pointer and added omitempty
}

func NewResponse(code int, status bool, message string, data interface{}) APIBaseResponse {
	return APIBaseResponse{
		Code:       code,
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: nil,
	}
}

func (response APIBaseResponse) NewResponseWithPagination(page int, limit int, total int64) APIBaseResponse {
	totalPage := (limit + int(total) - 1) / limit
	response.Pagination = &Pagination{
		Page:      page,
		Limit:     limit,
		TotalPage: totalPage,
		TotalData: total,
	}
	return response
}
