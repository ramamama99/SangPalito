package response

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"` 
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func OK(data interface{}) SuccessResponse {
	return SuccessResponse{Success: true, Data: data}
}

func OKWithMeta(data interface{}, meta Meta) SuccessResponse {
	return SuccessResponse{Success: true, Data: data, Meta: &meta}
}

func Error(code, message string) ErrorResponse {
	return ErrorResponse{Success: false, Code: code, Message: message}
}