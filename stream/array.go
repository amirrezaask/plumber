package stream

import (
	"github.com/amirrezaask/plumber"
	"io"
)

//ArrayInput uses an array to feed data into pipeline.
type ArrayInput struct {
	arr      []interface{}
	readChan chan interface{}
}

func NewArrayInput(words ...interface{}) plumber.Input {
	a := &ArrayInput{arr: words, readChan: make(chan interface{})}
	go func() {
		for _, e := range a.arr {
			a.readChan <- e
		}
	}()
	return a
}

func (a *ArrayInput) LoadState(r io.Reader) error {
	return nil
}

func (a *ArrayInput) Input() (chan interface{}, error) {
	return a.readChan, nil
}

func (a *ArrayInput) State() ([]byte, error) {
	return nil, nil
}

func (a *ArrayInput) Name() string {
	return "array-stream"
}
