package models

type TectonicError struct {
	status  int
	ErrCode uint   `json:"code"`
	Msg     string `json:"message"`
	Det     any    `json:"details,omitempty"`
}

func (e *TectonicError) Error() string  { return e.Msg }
func (e *TectonicError) GetStatus() int { return e.status }

func NewTectonicError(apiErr APIV1Error) *TectonicError {
	return &TectonicError{
		status:  apiErr.Status(),
		ErrCode: apiErr.Code(),
		Msg:     apiErr.Message(),
	}
}

// For when you need to attach extra context (e.g. field-level details)
func NewTectonicErrorWithDetails(apiErr APIV1Error, details any) *TectonicError {
	return &TectonicError{
		status:  apiErr.Status(),
		ErrCode: apiErr.Code(),
		Msg:     apiErr.Message(),
		Det:     details,
	}
}
