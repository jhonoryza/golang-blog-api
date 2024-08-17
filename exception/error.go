package exception

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{
		Error: error,
	}
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicNotFoundIfErr(err error) {
	if err != nil {
		panic(NewNotFoundError(err.Error()))
	}
}
