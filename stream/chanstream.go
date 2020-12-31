package stream

import (
	"errors"
	"github.com/amirrezaask/plumber"
	"io"
)

type ChanInput struct {
	readChan chan interface{}
	c        chan interface{}
}

func NewChanInput() plumber.Input {
	st := &ChanInput{c: make(chan interface{}), readChan: make(chan interface{})}
	go func() {
		for v := range st.c {
			st.readChan <- v
		}
	}()

	return st
}

func (d *ChanInput) Name() string {
	return "chan-stream"
}

func (d *ChanInput) LoadState(r io.Reader) error {
	return errors.New("chan stream is stateless")
}

func (d *ChanInput) State() ([]byte, error) {
	return nil, nil
}

func (d *ChanInput) Input() (chan interface{}, error) {
	return d.readChan, nil
}
