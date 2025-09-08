package exception

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}
