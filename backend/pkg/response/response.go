package response

// SuccessResponse is the standard JSON envelope for successful API responses.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorResponse is the standard JSON envelope for error API responses.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains pagination metadata.
type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// OK wraps data in a success response.
func OK(data interface{}) SuccessResponse {
	return SuccessResponse{Success: true, Data: data}
}

// OKWithMeta wraps data with pagination metadata in a success response.
func OKWithMeta(data interface{}, meta Meta) SuccessResponse {
	return SuccessResponse{Success: true, Data: data, Meta: &meta}
}

// Error constructs an error response with the given code and message.
func Error(code, message string) ErrorResponse {
	return ErrorResponse{Success: false, Code: code, Message: message}
}
