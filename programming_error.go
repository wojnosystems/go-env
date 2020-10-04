package env_parser

// ErrProgramming is an error that is encountered due to the developer misusing the Env.Unmarshall method.
type ErrProgramming struct {
	msg string
}

func NewErrProgramming(msg string) *ErrProgramming {
	return &ErrProgramming{
		msg: msg,
	}
}

func (p ErrProgramming) Error() string {
	return "programming error: " + p.msg
}
