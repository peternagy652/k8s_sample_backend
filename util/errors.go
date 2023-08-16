package util

// Not found error
type NotFoundError struct{}

var NotFoundErrorString = "Could not find matching person"

func (m *NotFoundError) Error() string {
	return NotFoundErrorString
}

// Not implemented error
type NotImplementedError struct{}

var NotImplementedErrorString = "Feature not implemented"

func (i *NotImplementedError) Error() string {
	return NotImplementedErrorString
}
