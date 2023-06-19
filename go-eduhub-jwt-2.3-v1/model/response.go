package model

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewSuccessResponse(msg string) SuccessResponse {
	return SuccessResponse{
		Message: msg,
	}
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}
