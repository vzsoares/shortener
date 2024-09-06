package utils

type Body struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Code    CODES  `json:"code"`
}

func NewBody(data any, msg string, code CODES) *Body {
	return &Body{
		Data:    data,
		Message: msg,
		Code:    code,
	}
}
