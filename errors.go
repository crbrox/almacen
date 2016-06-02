package almacen

type Error struct {
	statusCode int
	message    string
}

var (
	ErrExisting         = &Error{statusCode: 400, message: "existing id"}
	ErrNotFound         = &Error{statusCode: 404, message: "not found"}
	ErrTooMany          = &Error{statusCode: 400, message: "too many"}
	ErrObjectExpected   = &Error{statusCode: 400, message: "expected object"}
	ErrIdNotString      = &Error{statusCode: 500, message: "ID is not a string"}
	ErrTraversingObject = &Error{statusCode: 400, message: "traversing object"}
)

func (e *Error) Error() string {
	return e.message
}
