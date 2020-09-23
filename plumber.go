package plumber

import (
	"errors"
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

type DumbState map[string]interface{}

func (d DumbState) Get(k string) (interface{}, error) {
	v := d[k]
	return v, nil
}
func (d DumbState) Set(k string, v interface{}) error {
	d[k] = v
	return nil
}

type DumbStream struct {
	Stream chan []byte
}

func (d *DumbStream) Read(p []byte) (n int, err error) {
	v := <-d.Stream
	if len(v) > len(p) {
		return -1, errors.New("size don't match")
	}
	p = []byte(v)
	return len(p), nil
}
func (d *DumbStream) Write(p []byte) (n int, err error) {
	d.Stream <- p
	return len(p), nil
}
func (d *DumbStream) Close() error {
	close(d.Stream)
	return nil
}
