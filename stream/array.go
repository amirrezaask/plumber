package stream

import "github.com/amirrezaask/plumber"

type ArrayStream struct {
	arr      []interface{}
	readChan chan interface{}
}

func NewArrayStream(words ...interface{}) (plumber.Stream, error) {
	return &ArrayStream{arr: words, readChan: make(chan interface{})}, nil
}

func (a *ArrayStream) LoadState(map[string]interface{}) {
	return
}
func (a *ArrayStream) Write(v interface{}) error {
	return nil
}
func (a *ArrayStream) StartReading() error {
	go func() {
		for _, e := range a.arr {
			a.readChan <- e
		}
	}()
	return nil
}
func (a *ArrayStream) ReadChan() chan interface{} {
	return a.readChan
}

func (a *ArrayStream) State() map[string]interface{} {
	return nil
}

func (a *ArrayStream) Name() string {
	return "array-stream"
}
