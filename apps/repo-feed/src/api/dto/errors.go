package dto

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}
