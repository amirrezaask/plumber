package stream

import "github.com/amirrezaask/plumber"

type ChanStream struct {
	c        chan interface{}
	readChan chan interface{}
}

func NewChanStream() plumber.Stream {
	return &ChanStream{c: make(chan interface{}), readChan: make(chan interface{})}
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
		d.readChan <- d.c
	}()
	return nil
}
