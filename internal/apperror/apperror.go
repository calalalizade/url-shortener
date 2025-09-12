package apperror

type Type string

const (
	Validation Type = "VALIDATION"
	NotFound   Type = "NOT_FOUND"
	Internal   Type = "INTERNAL"
)

type AppError struct {
	Type    Type
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
