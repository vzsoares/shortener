package tools

import "fmt"

type AppError struct {
	Message string
	Code    CODES
	Who     *string
}

type CODES string

const (
	CODE_DB_ITEM_NOT_FOUND CODES = "DBI404"
)

var (
	ItemNotFoundError = New("Item not found", CODE_DB_ITEM_NOT_FOUND)
)

func (e AppError) Error() string {
	return fmt.Sprintf("App Error:  code: %x, msg: %x, Who: %x",
		string(e.Message), string(e.Code), string(*e.Who))
}

func New(msg string, code CODES) *AppError {
	return &AppError{Message: msg, Code: code}
}
