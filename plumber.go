package plumber

import (
	"io"
)

type Stream io.ReadWriteCloser

//Each state backend should implement this.
type State interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

// Lambda is a stateful function
type Lambda func(state State, in Stream, out Stream, errs Stream)

// Runner runs functions
type Runner interface {
	SetLogger(interface{}) //TODO: define logger interface
	Start(name string, fn Lambda) chan error
	Stop(name string) error
}
