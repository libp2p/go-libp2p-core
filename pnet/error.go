package pnet

// Error is error type for ease of detecting PNet errors
type Error interface {
	IsPNetError() bool
}

// NewError creates new Error
func NewError(err string) error {
	return pnetErr("privnet: " + err)
}

// IsPNetError checks if given error is PNet Error
func IsPNetError(err error) bool {
	v, ok := err.(Error)
	return ok && v.IsPNetError()
}

type pnetErr string

var _ Error = (*pnetErr)(nil)

func (p pnetErr) Error() string {
	return string(p)
}

func (pnetErr) IsPNetError() bool {
	return true
}
