package utils

type Error struct {
	Message string
	Cause error
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) getCause() error {
	return e.Cause
}
