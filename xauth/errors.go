package xauth

type Error struct {
	s string
}

func (e *Error) Error() string {
	return e.s
}

var (
	ErrMissingGroupMembership = &Error{"User not a member of one of the required groups"}
)
