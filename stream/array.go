package stream

import "github.com/amirrezaask/plumber"

type ArrayStream struct {
	arr      []interface{}
	readChan chan interface{}
}

func NewArrayStream(words ...interface{}) (plumber.Stream, error) {
	a := &ArrayStream{arr: words, readChan: make(chan interface{})}
	go func() {
		for _, e := range a.arr {
			a.readChan <- e
		}
	}()
	return a, nil
}

func (a *ArrayStream) LoadState(map[string]interface{}) {
	return
}
func (a *ArrayStream) Output() chan interface{} {
	panic("Array stream should only be used for intput")
}
func (a *ArrayStream) Input() chan interface{} {
	return a.readChan
}

func (a *ArrayStream) State() map[string]interface{} {
	return nil
}

func (a *ArrayStream) Name() string {
	return "array-stream"
}
