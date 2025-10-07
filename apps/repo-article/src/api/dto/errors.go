package dto

type ErrConflict struct{ Msg string }

func (e *ErrConflict) Error() string { return e.Msg }
func NewConflict(msg string) error   { return &ErrConflict{Msg: msg} }

type ErrNotFound struct{ Msg string }

func (e *ErrNotFound) Error() string { return e.Msg }
func NewNotFound(msg string) error   { return &ErrNotFound{Msg: msg} }

type ErrBadRequest struct{ Msg string }

func (e *ErrBadRequest) Error() string { return e.Msg }
func NewBadRequest(msg string) error   { return &ErrBadRequest{Msg: msg} }

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
