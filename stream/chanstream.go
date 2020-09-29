package stream

import "github.com/amirrezaask/plumber"

type ChanStream struct {
	c        chan interface{}
	readChan chan interface{}
}

//TODO: update according to stream constructor
func NewChanStream() plumber.Stream {
	st := &ChanStream{c: make(chan interface{}), readChan: make(chan interface{})}
	st.StartReading()
	return st
}

func (d *ChanStream) Name() string {
	return "chan-stream"
}

func (d *ChanStream) LoadState(m map[string]interface{}) {
	return
}
func (d *ChanStream) State() map[string]interface{} {
	return map[string]interface{}{}
}
func (d *ChanStream) Write(v interface{}) error {
	d.c <- v
	return nil
}
func (d *ChanStream) ReadChan() chan interface{} {
	return d.readChan
}
func (d *ChanStream) StartReading() error {
	go func() {
		for {
			d.readChan <- <-d.c
		}
	}()
	return nil
}
