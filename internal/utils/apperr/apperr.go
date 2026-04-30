package apperr

type AppErr struct {
	HTTPCode     int    `json:"-"`
	BusinessCode int    `json:"code"`
	Message      string `json:"message"`
}

func (e *AppErr) Error() string {
	return e.Message
}
