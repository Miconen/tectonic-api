package models

type TectonicError struct {
	status  int
	ErrCode uint   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *TectonicError) Error() string  { return e.Message }
func (e *TectonicError) GetStatus() int { return e.status }

func NewTectonicError(apiErr APIV1Error) *TectonicError {
	resp := apiErr.ToErrorResponse()
	return &TectonicError{
		status:  apiErr.Status(),
		ErrCode: resp.Code,
		Message: resp.Message,
		Details: resp.Details,
	}
}
