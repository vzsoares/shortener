package tools

import "fmt"

type CODES string

type AppError struct {
	Code CODES
}

func (e AppError) Error() string {
	return fmt.Sprintf("App Error:  code: %x", string(e.Code))
}

func New(code CODES) *AppError {
	return &AppError{Code: code}
}

const (
	CODE_BAD_REQUEST            CODES = "E400"
	CODE_OK                     CODES = "S200"
	CODE_INTERNAL_SERVER_ERROR  CODES = "E500"
	CODE_INPUT_VALIDATION_ERROR CODES = "E700"
	CODE_DB_ITEM_NOT_FOUND      CODES = "DBI404"
	CODE_DB_ITEM_ID_MISMATCH    CODES = "DBI40024"
)

var (
	ItemNotFoundError    = New(CODE_DB_ITEM_NOT_FOUND)
	IdMismatchError      = New(CODE_DB_ITEM_ID_MISMATCH)
	InputValidationError = New(CODE_INPUT_VALIDATION_ERROR)
)
