package response

const (
	StatusOk = "OK"
	StatusError = "Error"	
)

type Responce struct{
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}

func OK() Responce{
	return Responce{
		Status: StatusOk,
	}
}

func Error (msg string) Responce{
	return Responce{
		Status: StatusError,
		Error: msg,
	}
}